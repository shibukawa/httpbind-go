// Package htmlbind provides runtime configuration and response helpers for
// code generated from htmlbind templates.
package htmlbind

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/shibukawa/tinygodriver/compress/zstd"
)

// ZstdCompression controls whether generated HTML responses may use Zstandard
// content encoding. It is disabled by default. Set it during application
// startup, before serving requests; do not mutate it concurrently with requests.
var ZstdCompression bool

// PrepareResponse wraps w in a Zstandard encoder when compression is enabled
// and the request accepts zstd. The returned close function must always be
// called, including when rendering fails.
func PrepareResponse(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, func() error, error) {
	if !ZstdCompression {
		return w, noopClose, nil
	}
	addVary(w.Header(), "Accept-Encoding")
	if r == nil || !acceptsZstd(r.Header.Values("Accept-Encoding")) {
		return w, noopClose, nil
	}
	// An exported component may be rendered from another exported component.
	// In that case the outer response writer is already compressed.
	if w.Header().Get("Content-Encoding") != "" {
		return w, noopClose, nil
	}

	w.Header().Set("Content-Encoding", zstd.ContentEncoding)
	w.Header().Del("Content-Length")
	encoder, err := zstd.NewWriter(w, zstd.WithETag(false))
	if err != nil {
		return w, noopClose, err
	}
	return &encodedResponseWriter{ResponseWriter: w, writer: encoder}, encoder.Close, nil
}

func noopClose() error { return nil }

type encodedResponseWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (w *encodedResponseWriter) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

// Unwrap lets http.ResponseController reach optional interfaces implemented by
// the original response writer.
func (w *encodedResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func acceptsZstd(values []string) bool {
	for _, value := range values {
		for _, entry := range strings.Split(value, ",") {
			parts := strings.Split(entry, ";")
			if !strings.EqualFold(strings.TrimSpace(parts[0]), zstd.ContentEncoding) {
				continue
			}
			quality := 1.0
			for _, parameter := range parts[1:] {
				name, raw, ok := strings.Cut(parameter, "=")
				if !ok || !strings.EqualFold(strings.TrimSpace(name), "q") {
					continue
				}
				parsed, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
				if err != nil || parsed < 0 || parsed > 1 {
					quality = 0
				} else {
					quality = parsed
				}
			}
			if quality > 0 {
				return true
			}
		}
	}
	return false
}

func addVary(header http.Header, value string) {
	for _, line := range header.Values("Vary") {
		for _, existing := range strings.Split(line, ",") {
			existing = strings.TrimSpace(existing)
			if existing == "*" || strings.EqualFold(existing, value) {
				return
			}
		}
	}
	header.Add("Vary", value)
}
