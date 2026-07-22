package pages

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderedOutput(t *testing.T) {
	output := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	err := Document(
		output,
		request,
		"<b>raw</b>",
		"body > p { color: red; }",
		"window.ready = true;",
		Payload{Message: "<unsafe>&", Count: 2, Enabled: true},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := "\n<b>raw</b>\n<style>body > p { color: red; }</style>\n" +
		"<script>window.ready = true;</script>\n" +
		`<script>window.payload = {"message":"\u003cunsafe\u003e\u0026","count":2,"enabled":true};</script>` + "\n"
	if output.Body.String() != want {
		t.Fatalf("output = %q, want %q", output.Body.String(), want)
	}
}
