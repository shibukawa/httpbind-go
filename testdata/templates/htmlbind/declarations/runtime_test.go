package pages

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Decorate(value string, tone Tone) string {
	return string(tone) + ":" + value
}

func TestRenderedOutput(t *testing.T) {
	output := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	if err := Label(output, request, "<value>", TonePrimary); err != nil {
		t.Fatal(err)
	}
	want := "\n<span>Primary:&lt;value&gt;</span>\n"
	if output.Body.String() != want {
		t.Fatalf("output = %q, want %q", output.Body.String(), want)
	}
}
