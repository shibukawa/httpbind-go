package pages

import (
	"strings"
	"testing"
)

func TestRenderedOutput(t *testing.T) {
	var output strings.Builder
	if err := Hello(&output); err != nil {
		t.Fatal(err)
	}
	want := "\n<!DOCTYPE html>\n<h1>Hello &amp; welcome</h1>\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}
