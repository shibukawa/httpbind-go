package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ParsePackage analyzes Go sources in dir (same package only) and returns
// statically discoverable httpbinder route IR.
func ParsePackage(dir string) (*Result, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		if !strings.HasSuffix(name, ".go") {
			return false
		}
		// skip tests in the package under analysis
		if strings.HasSuffix(name, "_test.go") {
			return false
		}
		// skip generated binders / openapi embeds
		if strings.HasSuffix(name, "_httpbinder_gen.go") ||
			strings.HasSuffix(name, "_openapi_gen.go") ||
			name == "httpbinder_gen.go" ||
			name == "httpbinder_openapi_gen.go" {
			return false
		}
		return true
	}, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse dir %s: %w", dir, err)
	}
	if len(pkgs) == 0 {
		return &Result{Routes: []Route{}}, nil
	}

	// Same-package scope: if multiple packages appear (e.g. main + ignored),
	// prefer non-test package name; ParseDir already skipped _test.go.
	var pkg *ast.Package
	for name, p := range pkgs {
		if strings.HasSuffix(name, "_test") {
			continue
		}
		pkg = p
		break
	}
	if pkg == nil {
		for _, p := range pkgs {
			pkg = p
			break
		}
	}

	p := &packageParser{
		fset:  fset,
		pkg:   pkg,
		files: orderedFiles(pkg),
		funcs: map[string]*ast.FuncDecl{},
		types: map[string]*ast.TypeSpec{},
	}
	p.indexDecls()
	routes := p.discoverRoutes()
	res := &Result{Routes: routes}
	res.Normalize()
	return res, nil
}

// ParsePackageFiles is like ParsePackage but accepts an explicit file list
// (used by tests when embedding small snippets). Each path must exist.
func ParsePackageFiles(files []string) (*Result, error) {
	if len(files) == 0 {
		return &Result{Routes: []Route{}}, nil
	}
	dir := filepath.Dir(files[0])
	return ParsePackage(dir)
}

type packageParser struct {
	fset  *token.FileSet
	pkg   *ast.Package
	files []*ast.File
	funcs map[string]*ast.FuncDecl // name -> func (non-method)
	types map[string]*ast.TypeSpec
}

func orderedFiles(pkg *ast.Package) []*ast.File {
	names := make([]string, 0, len(pkg.Files))
	for name := range pkg.Files {
		names = append(names, name)
	}
	sortStrings(names)
	out := make([]*ast.File, 0, len(names))
	for _, name := range names {
		out = append(out, pkg.Files[name])
	}
	return out
}

func (p *packageParser) indexDecls() {
	for _, f := range p.files {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if fd.Recv == nil && fd.Name != nil {
				p.funcs[fd.Name.Name] = fd
			}
			// methods: Type.ServeHTTP tracked via type name elsewhere
		}
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if ok && ts.Name != nil {
					p.types[ts.Name.Name] = ts
				}
			}
		}
	}
}

func (p *packageParser) discoverRoutes() []Route {
	var routes []Route
	for _, f := range p.files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if route, ok := p.tryRouteCall(call); ok {
				routes = append(routes, route)
			}
			return true
		})
	}
	return routes
}

func (p *packageParser) tryRouteCall(call *ast.CallExpr) (Route, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel == nil {
		return Route{}, false
	}
	name := sel.Sel.Name
	if name != "HandleFunc" && name != "Handle" {
		return Route{}, false
	}
	// Receiver must be http package or a mux-like identifier (any x.Handle/HandleFunc).
	// We accept both http.HandleFunc and mux.HandleFunc.
	if len(call.Args) < 2 {
		return Route{}, false
	}

	patternLit, ok := stringLiteral(call.Args[0])
	if !ok {
		// Dynamic / non-static pattern → unsupported, no route discovery.
		return Route{}, false
	}
	method, path, ok := splitPattern(patternLit)
	if !ok {
		return Route{}, false
	}

	leaf, meta, ok := unwrapHandler(call.Args[1], WrapperMeta{})
	if !ok {
		return Route{}, false
	}

	handler, body := p.resolveHandler(leaf)
	if handler.Form == "" {
		return Route{}, false
	}

	route := Route{
		Method:   method,
		Path:     path,
		Handler:  handler,
		Wrappers: meta,
	}
	if body != nil {
		info := analyzeBody(body)
		route.Request = info.Request
		route.Response = info.Response
		route.Stream = info.Stream
		route.Errors = info.Errors
	}
	return route, true
}

func splitPattern(pattern string) (method, path string, ok bool) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return "", "", false
	}
	// Go 1.22+ "METHOD /path" or path-only.
	if i := strings.IndexByte(pattern, ' '); i >= 0 {
		m := strings.TrimSpace(pattern[:i])
		p := strings.TrimSpace(pattern[i+1:])
		if m == "" || p == "" {
			return "", "", false
		}
		return m, p, true
	}
	// path-only patterns are still valid registrations
	return "", pattern, true
}

func stringLiteral(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	s, err := strconvUnquote(lit.Value)
	if err != nil {
		return "", false
	}
	return s, true
}

func strconvUnquote(s string) (string, error) {
	// reuse strconv via thin wrapper to keep imports local in helpers
	return unquote(s)
}
