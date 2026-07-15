package httpbinder_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRuntimeMappingSources_NoFieldReflect walks shipped runtime mapping sources
// and ensures application field discovery is not done via the reflect package
// except the documented registry type-key helper in registry.go.
func TestRuntimeMappingSources_NoFieldReflect(t *testing.T) {
	root := "."
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(root, e.Name())
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		// registry.go may import reflect only for type identity keys.
		if e.Name() == "registry.go" {
			// type identity only: TypeFor (preferred) or TypeOf
			if !strings.Contains(string(src), "reflect.TypeFor") && !strings.Contains(string(src), "reflect.TypeOf") {
				t.Fatalf("registry.go should document type-key reflect usage")
			}
			// ensure no StructField / FieldByName style field walking
			for _, bad := range []string{"FieldByName", "NumField", "StructField", "Tag.Get"} {
				if strings.Contains(string(src), bad) {
					t.Fatalf("%s must not field-walk via reflect (%s)", path, bad)
				}
			}
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, src, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, imp := range f.Imports {
			if imp.Path != nil && imp.Path.Value == `"reflect"` {
				t.Fatalf("%s imports reflect; only registry.go may for type keys", path)
			}
		}
	}
}
