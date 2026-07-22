package pages

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	runtimehtmlbind "github.com/shibukawa/tinybind-go/htmlbind"
)

func TestRenderedOutput(t *testing.T) {
	output := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	if err := Hello(output, request); err != nil {
		t.Fatal(err)
	}
	if got := output.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	want := "\n<!DOCTYPE html>\n<h1>Hello &amp; welcome</h1>\n"
	if output.Body.String() != want {
		t.Fatalf("output = %q, want %q", output.Body.String(), want)
	}
}

func TestRenderedOutputZstd(t *testing.T) {
	runtimehtmlbind.ZstdCompression = true
	t.Cleanup(func() { runtimehtmlbind.ZstdCompression = false })
	output := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept-Encoding", "gzip, zstd")
	if err := Hello(output, request); err != nil {
		t.Fatal(err)
	}
	if got := output.Header().Get("Content-Encoding"); got != "zstd" {
		t.Fatalf("Content-Encoding = %q", got)
	}
	if !bytes.HasPrefix(output.Body.Bytes(), []byte{0x28, 0xb5, 0x2f, 0xfd}) {
		t.Fatalf("body is not a zstd frame: %x", output.Body.Bytes())
	}
}
