package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawayoshiki/httpbind-go/generator"
)

func TestAnalyzeAndEmit_NoReflect(t *testing.T) {
	dir := t.TempDir()
	src := `package sample

type Req struct {
	Name  string
	Page  int    ` + "`query:\"page\"`" + `
	OrgID string ` + "`path:\"org_id\"`" + `
	Token string ` + "`header:\"Authorization\"`" + `
	Note  string ` + "`payload:\"note\"`" + `
}

type Resp struct {
	ID   string ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`
	if err := os.WriteFile(filepath.Join(dir, "types.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	plan, err := generator.AnalyzePackage(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Types) != 2 {
		t.Fatalf("types: %+v", plan.Types)
	}
	code, err := generator.Emit(plan)
	if err != nil {
		t.Fatal(err)
	}
	s := string(code)
	if strings.Contains(s, "reflect") {
		t.Fatalf("reflect in output:\n%s", s)
	}
	for _, n := range []string{
		"func bindReq",
		"func writeResp",
		`QueryValue(r, "page")`,
		`PathValue(r, "org_id")`,
		`HeaderValue(r, "Authorization")`,
		`"note"`,
	} {
		if !strings.Contains(s, n) {
			t.Fatalf("missing %q in:\n%s", n, s)
		}
	}

	path, err := generator.Generate(dir, dir, "out_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(path)
	if err != nil || st.Size() == 0 {
		t.Fatalf("generated file empty or missing: %v", err)
	}
}
