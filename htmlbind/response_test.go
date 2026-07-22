package htmlbind

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shibukawa/tinygodriver/compress/zstd"
)

func TestPrepareResponseDefaultsToUncompressed(t *testing.T) {
	ZstdCompression = false
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept-Encoding", "zstd")

	w, closeResponse, err := PrepareResponse(recorder, request)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	if err := closeResponse(); err != nil {
		t.Fatal(err)
	}
	if got := recorder.Header().Get("Content-Encoding"); got != "" {
		t.Fatalf("Content-Encoding = %q", got)
	}
	if got := recorder.Body.String(); got != "hello" {
		t.Fatalf("body = %q", got)
	}
}

func TestPrepareResponseCompressesAcceptedZstd(t *testing.T) {
	ZstdCompression = true
	t.Cleanup(func() { ZstdCompression = false })
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Length", "5")
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept-Encoding", "gzip, zstd")

	w, closeResponse, err := PrepareResponse(recorder, request)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	if err := closeResponse(); err != nil {
		t.Fatal(err)
	}

	want, _, err := zstd.EncodeAll([]byte("hello"), zstd.WithETag(false))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(recorder.Body.Bytes(), want) {
		t.Fatalf("encoded body differs: got %x, want %x", recorder.Body.Bytes(), want)
	}
	if got := recorder.Header().Get("Content-Encoding"); got != "zstd" {
		t.Fatalf("Content-Encoding = %q", got)
	}
	if got := recorder.Header().Get("Content-Length"); got != "" {
		t.Fatalf("Content-Length = %q", got)
	}
	if got := recorder.Header().Values("Vary"); len(got) != 1 || got[0] != "Accept-Encoding" {
		t.Fatalf("Vary = %q", got)
	}
}

func TestPrepareResponseNegotiation(t *testing.T) {
	ZstdCompression = true
	t.Cleanup(func() { ZstdCompression = false })
	for _, value := range []string{"", "gzip", "zstd;q=0", "zstd;q=invalid"} {
		t.Run(value, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Accept-Encoding", value)
			w, closeResponse, err := PrepareResponse(recorder, request)
			if err != nil {
				t.Fatal(err)
			}
			if w != recorder {
				t.Fatalf("response was compressed for %q", value)
			}
			if err := closeResponse(); err != nil {
				t.Fatal(err)
			}
			if got := recorder.Header().Get("Vary"); got != "Accept-Encoding" {
				t.Fatalf("Vary = %q", got)
			}
		})
	}
}
