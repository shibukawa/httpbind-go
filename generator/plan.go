package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// FieldSource is where a request field is read from.
type FieldSource string

const (
	SourceInput   FieldSource = "input"
	SourceQuery   FieldSource = "query"
	SourcePayload FieldSource = "payload"
	SourcePath    FieldSource = "path"
	SourceHeader  FieldSource = "header"
	SourceCookie  FieldSource = "cookie"
	SourceMethod  FieldSource = "method"
)

// Composite and special field kinds.
const (
	KindRestAny = "rest_any" // map[string]any with payload:"*"
	KindRestRaw = "rest_raw" // map[string]json.RawMessage with payload:"*"
	KindStruct  = "struct"
	KindSlice   = "slice"
	KindMap     = "map"
)

// FieldPlan is one struct field mapping plan (compile-time).
type FieldPlan struct {
	Name     string      // Go field name
	Wire     string      // wire / tag name ("*" for payload rest)
	Source   FieldSource // input|query|payload|path|header|cookie|method
	Kind     string      // string|int|int64|bool|float64|file|rest_*|struct|slice|map
	JSON     string      // json name for encode/document keys
	Check    CheckRules  // from check:"" tag; empty if absent
	TypeName string      // KindStruct name, or element struct name for slice/map of struct
	ElemKind string      // for slice/map: string|int|int64|bool|float64|struct
}

// IsRest reports whether f is a payload rest map field.
func (f FieldPlan) IsRest() bool {
	return f.Kind == KindRestAny || f.Kind == KindRestRaw
}

// IsComposite reports nested struct/slice/map kinds.
func (f FieldPlan) IsComposite() bool {
	return f.Kind == KindStruct || f.Kind == KindSlice || f.Kind == KindMap
}

// GoType returns a Go type string for generated code (e.g. NestedCustomer, []string).
func (f FieldPlan) GoType() string {
	switch f.Kind {
	case KindStruct:
		return f.TypeName
	case KindSlice:
		if f.ElemKind == KindStruct {
			return "[]" + f.TypeName
		}
		return "[]" + f.ElemKind
	case KindMap:
		if f.ElemKind == KindStruct {
			return "map[string]" + f.TypeName
		}
		return "map[string]" + f.ElemKind
	case KindRestAny:
		return "map[string]any"
	case KindRestRaw:
		return "map[string]json.RawMessage"
	case "file":
		return "httpbinder.File"
	default:
		return f.Kind
	}
}

// httpbinderImportPath is the module path of this library (for recognizing File).
const httpbinderImportPath = "github.com/shibukawa/httpbind-go"

// TypePlan is the mapping plan for one struct type.
type TypePlan struct {
	Name   string
	Fields []FieldPlan
}

// PackagePlan is all type plans in a package.
type PackagePlan struct {
	Package string
	Types   []TypePlan
	// Discovered lists type names referenced by Bind/Write/DecodeJSON/EncodeJSON call sites.
	Discovered []string
}

// AnalyzePackage builds field plans for all package-level structs with exported fields.
func AnalyzePackage(dir string) (*PackagePlan, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return false
		}
		// skip generated files when re-analyzing
		if strings.HasSuffix(name, "_httpbinder_gen.go") ||
			strings.HasSuffix(name, "_openapi_gen.go") ||
			name == "httpbinder_gen.go" ||
			name == "httpbinder_openapi_gen.go" {
			return false
		}
		return true
	}, 0)
	if err != nil {
		return nil, err
	}
	var pkg *ast.Package
	for name, p := range pkgs {
		if strings.HasSuffix(name, "_test") {
			continue
		}
		pkg = p
		break
	}
	if pkg == nil {
		return nil, fmt.Errorf("no package in %s", dir)
	}

	plan := &PackagePlan{Package: pkg.Name}
	discovered := map[string]struct{}{}
	for _, f := range pkg.Files {
		binderNames := httpbinderImportNames(f)
		for _, name := range discoverGenericTypeArgs(f, binderNames) {
			discovered[name] = struct{}{}
		}
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name == nil {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok || st.Fields == nil {
					continue
				}
				tp, ok, err := analyzeStruct(ts.Name.Name, st, binderNames)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", ts.Name.Name, err)
				}
				if ok {
					plan.Types = append(plan.Types, tp)
				}
			}
		}
	}
	for name := range discovered {
		plan.Discovered = append(plan.Discovered, name)
	}
	// Ensure discovered types exist as plans (all exported structs are planned already).
	have := map[string]bool{}
	for _, t := range plan.Types {
		have[t.Name] = true
	}
	for name := range discovered {
		if !have[name] {
			// Type may be referenced but not a planned struct (e.g. missing exported fields).
			// Leave a clear error at generate time only if codecs are required; analysis allows it.
			_ = name
		}
	}
	return plan, nil
}

// discoverGenericTypeArgs finds type arguments of httpbinder.Bind/Write/DecodeJSON/EncodeJSON.
func discoverGenericTypeArgs(f *ast.File, binderNames map[string]bool) []string {
	var out []string
	if f == nil {
		return out
	}
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		name, typeArgs := genericCallInfo(call.Fun)
		if name == "" || len(typeArgs) == 0 {
			return true
		}
		switch name {
		case "Bind", "Write", "DecodeJSON", "EncodeJSON":
		default:
			return true
		}
		// Fun is pkg.Name[T] — ensure pkg is httpbinder
		if !callFunIsHTTPBinder(call.Fun, binderNames) {
			return true
		}
		for _, a := range typeArgs {
			if id, ok := a.(*ast.Ident); ok && id.Name != "" {
				out = append(out, id.Name)
			}
		}
		return true
	})
	return out
}

func callFunIsHTTPBinder(fun ast.Expr, binderNames map[string]bool) bool {
	// IndexExpr or IndexListExpr wrapping SelectorExpr
	switch f := fun.(type) {
	case *ast.IndexExpr:
		return selectorIsHTTPBinder(f.X, binderNames)
	case *ast.IndexListExpr:
		return selectorIsHTTPBinder(f.X, binderNames)
	case *ast.SelectorExpr:
		return selectorIsHTTPBinder(f, binderNames)
	}
	return false
}

func selectorIsHTTPBinder(expr ast.Expr, binderNames map[string]bool) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return binderNames[pkg.Name]
}

func genericCallInfo(fun ast.Expr) (name string, typeArgs []ast.Expr) {
	switch f := fun.(type) {
	case *ast.IndexExpr:
		if sel, ok := f.X.(*ast.SelectorExpr); ok && sel.Sel != nil {
			return sel.Sel.Name, []ast.Expr{f.Index}
		}
	case *ast.IndexListExpr:
		if sel, ok := f.X.(*ast.SelectorExpr); ok && sel.Sel != nil {
			return sel.Sel.Name, f.Indices
		}
	}
	return "", nil
}

// httpbinderImportNames returns local identifiers that refer to this library
// (default name "httpbinder" or explicit/aliased imports).
func httpbinderImportNames(f *ast.File) map[string]bool {
	out := make(map[string]bool)
	if f == nil {
		return out
	}
	for _, imp := range f.Imports {
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil || path != httpbinderImportPath {
			continue
		}
		if imp.Name != nil {
			switch imp.Name.Name {
			case "_":
				// side-effect import only
			case ".":
				// dot-import
			default:
				out[imp.Name.Name] = true
			}
			continue
		}
		out["httpbinder"] = true
	}
	return out
}

func analyzeStruct(name string, st *ast.StructType, binderNames map[string]bool) (TypePlan, bool, error) {
	var fields []FieldPlan
	restCount := 0
	for _, f := range st.Fields.List {
		if len(f.Names) == 0 {
			continue // embedded
		}
		for _, id := range f.Names {
			if id == nil || !exported(id.Name) {
				continue
			}
			src, wire := parseFieldTag(id.Name, f.Tag)
			fp, ok, err := analyzeField(id.Name, f.Type, f.Tag, src, wire, binderNames)
			if err != nil {
				return TypePlan{}, false, err
			}
			if !ok {
				continue
			}
			if fp.Kind == "file" {
				switch src {
				case SourceInput, SourcePayload:
					fp.Source = SourcePayload
				default:
					return TypePlan{}, false, fmt.Errorf("field %s: httpbinder.File only supports payload/input tags, got %s", id.Name, src)
				}
			}
			if fp.IsRest() {
				fp.Source = SourcePayload
				fp.Wire = "*"
				restCount++
				if restCount > 1 {
					return TypePlan{}, false, fmt.Errorf("field %s: at most one payload:\"*\" rest field allowed", id.Name)
				}
			}
			// Nested composites are JSON-oriented; force payload when tagged input for body nesting.
			if fp.IsComposite() {
				switch fp.Source {
				case SourceInput, SourcePayload:
					// keep; JSON bind uses body
				case SourceQuery, SourcePath, SourceHeader, SourceCookie, SourceMethod:
					return TypePlan{}, false, fmt.Errorf("field %s: nested %s only supports payload/input sources", id.Name, fp.Kind)
				}
			}
			fields = append(fields, fp)
		}
	}
	if len(fields) == 0 {
		return TypePlan{}, false, nil
	}
	return TypePlan{Name: name, Fields: fields}, true, nil
}

func analyzeField(fieldName string, typ ast.Expr, tag *ast.BasicLit, src FieldSource, wire string, binderNames map[string]bool) (FieldPlan, bool, error) {
	kind, typeName, elemKind, ok, err := fieldTypeKind(typ, binderNames, src, wire, fieldName)
	if err != nil {
		return FieldPlan{}, false, err
	}
	if !ok {
		return FieldPlan{}, false, nil
	}
	jsonName := wire
	if jsonName == "" || jsonName == "*" {
		jsonName = lowerFirst(fieldName)
	}
	if jt := tagValue(tag, "json"); jt != "" && jt != "-" {
		jsonName = strings.Split(jt, ",")[0]
	}
	checkRaw := tagValue(tag, "check")
	check, err := ParseCheckTag(checkRaw, kind)
	if err != nil {
		return FieldPlan{}, false, fmt.Errorf("field %s: %w", fieldName, err)
	}
	return FieldPlan{
		Name:     fieldName,
		Wire:     wire,
		Source:   src,
		Kind:     kind,
		JSON:     jsonName,
		Check:    check,
		TypeName: typeName,
		ElemKind: elemKind,
	}, true, nil
}

func exported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

// fieldTypeKind resolves a field's bind kind.
func fieldTypeKind(expr ast.Expr, binderNames map[string]bool, src FieldSource, wire, fieldName string) (kind, typeName, elemKind string, ok bool, err error) {
	if restKind, isRest := mapRestKind(expr); isRest {
		if wire != "*" {
			return "", "", "", false, nil
		}
		if src != SourcePayload {
			return "", "", "", false, fmt.Errorf("field %s: rest map requires payload:\"*\", got %s:%q", fieldName, src, wire)
		}
		return restKind, "", "", true, nil
	}
	if wire == "*" {
		return "", "", "", false, fmt.Errorf("field %s: payload:\"*\" requires map[string]any or map[string]json.RawMessage", fieldName)
	}

	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "string", "int", "int64", "bool", "float64":
			return t.Name, "", "", true, nil
		case "any", "error":
			return "", "", "", false, nil
		default:
			// Named type in the same package → nested struct.
			if t.Name != "" {
				return KindStruct, t.Name, "", true, nil
			}
		}
	case *ast.SelectorExpr:
		if t.Sel != nil && t.Sel.Name == "File" {
			if pkg, ok := t.X.(*ast.Ident); ok && binderNames[pkg.Name] {
				return "file", "", "", true, nil
			}
		}
	case *ast.ArrayType:
		ek, et, _, eok, eerr := fieldTypeKind(t.Elt, binderNames, src, wire, fieldName)
		if eerr != nil {
			return "", "", "", false, eerr
		}
		if !eok {
			return "", "", "", false, nil
		}
		switch ek {
		case "string", "int", "int64", "bool", "float64":
			return KindSlice, "", ek, true, nil
		case KindStruct:
			return KindSlice, et, KindStruct, true, nil
		default:
			return "", "", "", false, nil
		}
	case *ast.MapType:
		key, ok := t.Key.(*ast.Ident)
		if !ok || key.Name != "string" {
			return "", "", "", false, nil
		}
		ek, et, _, eok, eerr := fieldTypeKind(t.Value, binderNames, src, wire, fieldName)
		if eerr != nil {
			return "", "", "", false, eerr
		}
		if !eok {
			return "", "", "", false, nil
		}
		switch ek {
		case "string", "int", "int64", "bool", "float64":
			return KindMap, "", ek, true, nil
		case KindStruct:
			return KindMap, et, KindStruct, true, nil
		default:
			return "", "", "", false, nil
		}
	}
	return "", "", "", false, nil
}

func mapRestKind(expr ast.Expr) (string, bool) {
	mt, ok := expr.(*ast.MapType)
	if !ok {
		return "", false
	}
	key, ok := mt.Key.(*ast.Ident)
	if !ok || key.Name != "string" {
		return "", false
	}
	switch v := mt.Value.(type) {
	case *ast.Ident:
		if v.Name == "any" {
			return KindRestAny, true
		}
	case *ast.InterfaceType:
		if v.Methods == nil || len(v.Methods.List) == 0 {
			return KindRestAny, true
		}
	case *ast.SelectorExpr:
		if v.Sel != nil && v.Sel.Name == "RawMessage" {
			if pkg, ok := v.X.(*ast.Ident); ok && pkg.Name == "json" {
				return KindRestRaw, true
			}
		}
	}
	return "", false
}

func parseFieldTag(fieldName string, tag *ast.BasicLit) (FieldSource, string) {
	defaultWire := lowerFirst(fieldName)
	if tag == nil {
		return SourceInput, defaultWire
	}
	raw, err := strconv.Unquote(tag.Value)
	if err != nil {
		return SourceInput, defaultWire
	}
	for _, src := range []FieldSource{SourceInput, SourceQuery, SourcePayload, SourcePath, SourceHeader, SourceCookie, SourceMethod} {
		if v := lookupTag(raw, string(src)); v != "" {
			if v == "-" {
				continue
			}
			return src, v
		}
	}
	return SourceInput, defaultWire
}

func tagValue(tag *ast.BasicLit, key string) string {
	if tag == nil {
		return ""
	}
	raw, err := strconv.Unquote(tag.Value)
	if err != nil {
		return ""
	}
	return lookupTag(raw, key)
}

func lookupTag(raw, key string) string {
	for _, part := range strings.Fields(raw) {
		k, v, ok := strings.Cut(part, ":")
		if !ok || k != key {
			continue
		}
		val, err := strconv.Unquote(v)
		if err != nil {
			return strings.Trim(v, `"`)
		}
		return val
	}
	return ""
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}
