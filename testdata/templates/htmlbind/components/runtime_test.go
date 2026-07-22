package pages

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderedOutput(t *testing.T) {
	output := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	if err := Card(output, request, User{Name: "A&B"}); err != nil {
		t.Fatal(err)
	}
	want := `<span class="badge"><strong>A&amp;B</strong><em>member</em></span>`
	if !strings.Contains(output.Body.String(), want) {
		t.Fatalf("output %q does not contain %q", output.Body.String(), want)
	}
}
