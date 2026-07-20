package htmlbind_test

import (
	"bytes"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/shibukawa/tinybind-go/templates/htmlbind"
)

func TestGenerateFixtures(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "templates", "htmlbind")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	var cases []string
	for _, entry := range entries {
		if entry.IsDir() {
			cases = append(cases, entry.Name())
		}
	}
	sort.Strings(cases)
	if len(cases) == 0 {
		t.Fatal("no HTML generator fixtures found")
	}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(root, name)
			inputPath := filepath.Join(dir, "input.txt")
			outputPath := filepath.Join(dir, "output.go")
			input, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatal(err)
			}
			got, err := htmlbind.Generate(inputPath, input, htmlbind.GenerateOptions{})
			if err != nil {
				t.Fatal(err)
			}
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				if err := os.WriteFile(outputPath, got, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("generated Go mismatch\n--- want ---\n%s--- got ---\n%s", want, got)
			}
			runtimeTest, err := os.ReadFile(filepath.Join(dir, "runtime_test.go"))
			if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			typeCheckGenerated(t, outputPath, got, runtimeTest)
			runGeneratedTests(t, got, runtimeTest)
		})
	}
}

func runGeneratedTests(t *testing.T, generated, runtimeTest []byte) {
	t.Helper()
	if len(runtimeTest) == 0 {
		return
	}
	dir := t.TempDir()
	for name, content := range map[string][]byte{
		"go.mod":          []byte("module generatedfixture\n\ngo 1.26\n"),
		"generated.go":    generated,
		"runtime_test.go": runtimeTest,
	} {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	command := exec.Command("go", "test", ".")
	command.Dir = dir
	command.Env = append(os.Environ(), "GOWORK=off")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("run generated Go: %v\n%s", err, output)
	}
}

func typeCheckGenerated(t *testing.T, filename string, source, companion []byte) {
	t.Helper()
	files := token.NewFileSet()
	file, err := parser.ParseFile(files, filename, source, parser.AllErrors)
	if err != nil {
		t.Fatalf("parse generated Go: %v", err)
	}
	parsed := []*ast.File{file}
	if len(companion) > 0 {
		companionFile, err := parser.ParseFile(files, "runtime_test.go", companion, parser.AllErrors)
		if err != nil {
			t.Fatalf("parse generated Go companion: %v", err)
		}
		parsed = append(parsed, companionFile)
	}
	config := types.Config{Importer: importer.Default(), Error: func(err error) { t.Errorf("generated Go type error: %v", err) }}
	if _, err := config.Check(file.Name.Name, files, parsed, nil); err != nil {
		t.Fatalf("type-check generated Go: %v", err)
	}
}

func TestGenerateDiagnostics(t *testing.T) {
	tests := []struct{ name, source, want string }{
		{"unknown identifier", `component Bad(): html {<p>{missing}</p>}`, "unknown identifier missing"},
		{"wrong condition", `component Bad(name: string): html {{if name}x{/if}}`, "if condition must be bool"},
		{"unsafe script", `component Bad(value: string): html {<script>{value}</script>}`, "html:script requires"},
		{"unsafe raw context", `component Bad(value: string): html {<p title={RawHTML(value)}>x</p>}`, "cannot insert trusted_html"},
		{"optional raw input", `component Bad(value: string?): html {{RawHTML(value)}}`, "RawHTML expects string"},
		{"url type", `component Bad(value: string): html {<a href={value}>x</a>}`, "requires url"},
		{"optional mixed attribute", `component Bad(value: string?): html {<p title="prefix {value}">x</p>}`, "optional expression must be the entire attribute"},
		{"unsafe json field", `type Payload { target: url } component Bad(value: Payload): html {<script>{JsonForScript(value)}</script>}`, "not statically serializable"},
		{"noncomparable values", `component Bad(left: string[], right: string[]): html {{if left == right}x{/if}}`, "values are not comparable"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := htmlbind.Generate("invalid.txt", []byte(test.source), htmlbind.GenerateOptions{Package: "invalid"})
			if err == nil || !bytes.Contains([]byte(err.Error()), []byte(test.want)) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestGenerateManglesGoKeywords(t *testing.T) {
	source := []byte(`package type
export component Keyword(type: string): html {<p>{type}</p>}`)
	generated, err := htmlbind.Generate("keywords.txt", source, htmlbind.GenerateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(generated, []byte("package _type")) || !bytes.Contains(generated, []byte("_type string")) {
		t.Fatalf("generated Go does not mangle keywords:\n%s", generated)
	}
	typeCheckGenerated(t, "keywords.go", generated, nil)
}

func TestGenerateDiagnosticIncludesPosition(t *testing.T) {
	source := []byte("component Bad(): html {\n<p>\n{missing}\n</p>\n}")
	_, err := htmlbind.Generate("position.txt", source, htmlbind.GenerateOptions{})
	if err == nil || !bytes.Contains([]byte(err.Error()), []byte("position.txt:3:2:")) {
		t.Fatalf("error = %v, want filename:line:col", err)
	}
}
