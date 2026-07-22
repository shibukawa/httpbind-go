package pages

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRenderedOutput(t *testing.T) {
	nickname := "A & B"
	profileURL, err := url.Parse("https://example.com/profile?q=a&lang=en")
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		Name:       "<Ada>",
		Active:     true,
		Nickname:   &nickname,
		ProfileURL: *profileURL,
		Tags:       []string{"go", "<html>"},
	}
	output := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	if err := Profile(output, request, user); err != nil {
		t.Fatal(err)
	}
	rendered := output.Body.String()
	for _, want := range []string{
		`title="A &amp; B"`,
		`href="https://example.com/profile?q=a&amp;lang=en"`,
		`&lt;Ada&gt;`,
		`<li data-index="0">go</li>`,
		`<li data-index="1">&lt;html&gt;</li>`,
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("output %q does not contain %q", rendered, want)
		}
	}
	if strings.Contains(rendered, " hidden") || strings.Contains(rendered, "inactive") {
		t.Fatalf("unexpected inactive output: %q", rendered)
	}
}
