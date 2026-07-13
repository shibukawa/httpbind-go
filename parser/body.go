package parser

import (
	"go/ast"
	"sort"
)

type bodyInfo struct {
	Request  string
	Response string
	Stream   string
	Errors   []string
}

var errorConstructors = map[string]struct{}{
	"BadRequest":   {},
	"Unauthorized": {},
	"Forbidden":    {},
	"NotFound":     {},
	"Conflict":     {},
	"Internal":     {},
	"Validation":   {},
}

func analyzeBody(body *ast.BlockStmt) bodyInfo {
	var info bodyInfo
	if body == nil {
		return info
	}
	errorSet := map[string]struct{}{}
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		// httpbinder.Bind[T](...)
		if name, typeArgs, ok := genericCall(call); ok && name == "Bind" && len(typeArgs) >= 1 {
			if info.Request == "" {
				info.Request = typeArgs[0]
			}
			return true
		}
		// httpbinder.Write[T](...)
		if name, typeArgs, ok := genericCall(call); ok && name == "Write" && len(typeArgs) >= 1 {
			resp, streamElem := parseResponseType(typeArgs[0])
			if info.Response == "" {
				info.Response = resp
			}
			if streamElem != "" && info.Stream == "" {
				info.Stream = streamElem
			}
			return true
		}
		// httpbinder.NewStream[T](...) — preferred incremental streaming API
		if name, typeArgs, ok := genericCall(call); ok && name == "NewStream" && len(typeArgs) >= 1 {
			elem := typeArgs[0]
			if info.Stream == "" {
				info.Stream = elem
			}
			if info.Response == "" {
				info.Response = "httpbinder.Stream[" + elem + "]"
			}
			return true
		}
		// httpbinder.BadRequest / etc.
		if name, ok := httpbinderSelector(call.Fun); ok {
			if _, known := errorConstructors[name]; known {
				errorSet[name] = struct{}{}
			}
		}
		return true
	})
	if len(errorSet) > 0 {
		info.Errors = make([]string, 0, len(errorSet))
		for k := range errorSet {
			info.Errors = append(info.Errors, k)
		}
		sort.Strings(info.Errors)
	}
	return info
}

// genericCall recognizes pkg.Name[T](args) or Name[T](args).
// Returns the function name and stringified type arguments.
func genericCall(call *ast.CallExpr) (name string, typeArgs []string, ok bool) {
	switch fun := call.Fun.(type) {
	case *ast.IndexExpr:
		// Bind[T]
		name, ok = selectorOrIdentName(fun.X)
		if !ok {
			return "", nil, false
		}
		return name, []string{typeExprString(fun.Index)}, true
	case *ast.IndexListExpr:
		// Bind[T, U]
		name, ok = selectorOrIdentName(fun.X)
		if !ok {
			return "", nil, false
		}
		args := make([]string, 0, len(fun.Indices))
		for _, idx := range fun.Indices {
			args = append(args, typeExprString(idx))
		}
		return name, args, true
	default:
		return "", nil, false
	}
}

func selectorOrIdentName(expr ast.Expr) (string, bool) {
	expr = stripParens(expr)
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name, true
	case *ast.SelectorExpr:
		if e.Sel != nil {
			// Prefer requiring httpbinder package when selector
			if id, ok := e.X.(*ast.Ident); ok {
				_ = id // package name may be aliased; accept any
			}
			return e.Sel.Name, true
		}
	}
	return "", false
}

func httpbinderSelector(fun ast.Expr) (string, bool) {
	fun = stripParens(fun)
	// Match only httpbinder.NotFound(...) etc. — not stdlib http.NotFound.
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok || sel.Sel == nil {
		return "", false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "httpbinder" {
		return "", false
	}
	return sel.Sel.Name, true
}

// parseResponseType turns Write type arg into response name + optional stream element.
// Handles: CreateUserResponse, Stream[ChatEvent], httpbinder.Stream[ChatEvent]
func parseResponseType(typeStr string) (response, streamElem string) {
	// typeExprString already produces compact forms like "Stream[ChatEvent]"
	// or "httpbinder.Stream[ChatEvent]"
	s := typeStr
	// strip package prefix for Stream
	if i := lastIndexStream(s); i >= 0 {
		// e.g. httpbinder.Stream[ChatEvent] or Stream[ChatEvent]
		inner := extractBracketContent(s[i:])
		if inner != "" {
			return s, inner
		}
	}
	return s, ""
}

func lastIndexStream(s string) int {
	// find Stream[
	for i := 0; i+7 <= len(s); i++ {
		if s[i:i+7] == "Stream[" {
			return i
		}
	}
	return -1
}

func extractBracketContent(s string) string {
	// s starts at Stream[...]
	start := -1
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '[':
			if depth == 0 {
				start = i + 1
			}
			depth++
		case ']':
			depth--
			if depth == 0 && start >= 0 {
				return s[start:i]
			}
		}
	}
	return ""
}

func typeExprString(expr ast.Expr) string {
	expr = stripParens(expr)
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		left := typeExprString(e.X)
		if e.Sel == nil {
			return left
		}
		if left == "" {
			return e.Sel.Name
		}
		return left + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + typeExprString(e.X)
	case *ast.IndexExpr:
		return typeExprString(e.X) + "[" + typeExprString(e.Index) + "]"
	case *ast.IndexListExpr:
		parts := make([]string, 0, len(e.Indices))
		for _, idx := range e.Indices {
			parts = append(parts, typeExprString(idx))
		}
		return typeExprString(e.X) + "[" + joinComma(parts) + "]"
	case *ast.ArrayType:
		if e.Len == nil {
			return "[]" + typeExprString(e.Elt)
		}
		return "[" + typeExprString(e.Len) + "]" + typeExprString(e.Elt)
	case *ast.MapType:
		return "map[" + typeExprString(e.Key) + "]" + typeExprString(e.Value)
	default:
		return ""
	}
}

func joinComma(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += ", " + parts[i]
	}
	return out
}
