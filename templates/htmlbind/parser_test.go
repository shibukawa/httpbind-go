package htmlbind_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/shibukawa/tinybind-go/templates/htmlbind"
)

func TestParserFixtures(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "templates", "htmlparser")
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
		t.Fatal("no parser fixtures found")
	}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(root, name)
			inputPath := filepath.Join(dir, "input.txt")
			astPath := filepath.Join(dir, "ast.json")
			input, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatal(err)
			}
			module, err := htmlbind.Parse(inputPath, input)
			if err != nil {
				t.Fatal(err)
			}
			got, err := json.MarshalIndent(module, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			got = append(got, '\n')
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				if err := os.WriteFile(astPath, got, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(astPath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("AST mismatch\n--- want ---\n%s--- got ---\n%s", want, got)
			}
		})
	}
}

func TestParserDiagnostics(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{"output mismatch", `component Bad(): string {<p>x</p>}`, "component declaration requires html output"},
		{"mismatched tag", `component Bad(): html {<div></span>}`, "expected closing tag </div>"},
		{"control in attribute", `component Bad(ok: bool): html {<p class="{if ok}x{/if}">x</p>}`, "control blocks are forbidden in attributes"},
		{"unknown root", `statement Bad(): sql.exec {SELECT 1}`, "unknown root declaration"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := htmlbind.Parse("invalid.txt", []byte(tt.source))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want substring %q", err, tt.want)
			}
			var parseErr *htmlbind.ParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("error type = %T, want *htmlbind.ParseError", err)
			}
			if parseErr.Line < 1 || parseErr.Column < 1 {
				t.Fatalf("invalid position: %+v", parseErr)
			}
		})
	}
}

func TestASTPositionsUseFileLinesAndRuneColumns(t *testing.T) {
	source := "component Card(name: string): html {\n<p title={name}>あ{name}</p>\n}"
	module, err := htmlbind.Parse("position.txt", []byte(source))
	if err != nil {
		t.Fatal(err)
	}
	decl := module.Declarations[0].(*htmlbind.TemplateDecl)
	if decl.Pos != (htmlbind.Position{Line: 1, Col: 1}) {
		t.Fatalf("declaration pos = %+v", decl.Pos)
	}
	body := decl.Body.(htmlbind.Body)
	element := body[1].(*htmlbind.ElementNode)
	if element.Pos != (htmlbind.Position{Line: 2, Col: 1}) {
		t.Fatalf("element pos = %+v", element.Pos)
	}
	if element.Attributes[0].Pos != (htmlbind.Position{Line: 2, Col: 4}) {
		t.Fatalf("attribute pos = %+v", element.Attributes[0].Pos)
	}
	attributeExpr := element.Attributes[0].Value[0]
	if attributeExpr.Pos != (htmlbind.Position{Line: 2, Col: 10}) {
		t.Fatalf("attribute expression pos = %+v", attributeExpr.Pos)
	}
	text := element.Children[0].(*htmlbind.TextNode)
	if text.Pos != (htmlbind.Position{Line: 2, Col: 17}) {
		t.Fatalf("unicode text pos = %+v", text.Pos)
	}
	expr := element.Children[1].(*htmlbind.ExpressionNode)
	if expr.Pos != (htmlbind.Position{Line: 2, Col: 18}) {
		t.Fatalf("expression pos after unicode = %+v", expr.Pos)
	}
}

func TestDiagnosticPositionUsesWholeFile(t *testing.T) {
	source := "package demo\n\ncomponent Bad(): html {\n<div></span>\n}"
	_, err := htmlbind.Parse("position.txt", []byte(source))
	var parseErr *htmlbind.ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("error = %v", err)
	}
	if parseErr.Line != 4 {
		t.Fatalf("diagnostic line = %d, want 4: %v", parseErr.Line, err)
	}
}
