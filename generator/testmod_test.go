package generator_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// writeTempModule writes a go.mod that replace-points at this module root so
// packages.Load can type-check temp packages that import tinybind-go.
func writeTempModule(t *testing.T, dir string) {
	t.Helper()
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	mod := "module tempmod\n\n" +
		"go 1.25\n\n" +
		"require github.com/shibukawa/tinybind-go v0.0.0\n\n" +
		"replace github.com/shibukawa/tinybind-go => " + filepath.ToSlash(root) + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o644); err != nil {
		t.Fatal(err)
	}
}

// tidyTempModule runs go mod tidy after package sources are written.
func tidyTempModule(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}
}
