package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/tinybind-go/generator"
)

func TestAnalyze_DiscoversAliasedDecodeEncode(t *testing.T) {
	dir := t.TempDir()
	writeTempModule(t, dir)
	src := `package sample

import (
	hb "github.com/shibukawa/tinybind-go/jsonbind"
)

type Note struct {
	Text string ` + "`payload:\"text\"`" + `
}

func use() {
	_, _ = hb.DecodeJSON[Note](nil)
	_ = hb.EncodeJSON(nil, Note{})
}

// Foreign same-named helpers must not contribute discovery.
func DecodeJSON[T any](r any) (T, error) {
	var zero T
	return zero, nil
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	tidyTempModule(t, dir)
	plan, err := generator.AnalyzePackage(dir)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, n := range plan.Discovered {
		if n == "Note" {
			found = true
		}
	}
	if !found {
		t.Fatalf("alias DecodeJSON discovery missing Note: %v", plan.Discovered)
	}
	code, err := generator.Emit(plan)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(code), "RegisterDecode[Note]") {
		t.Fatal("missing RegisterDecode[Note]")
	}
}
