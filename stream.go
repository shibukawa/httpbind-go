package httpbind

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// StreamFormat is the negotiated on-the-wire format for Stream[T].
type StreamFormat string

const (
	// StreamSSE is text/event-stream (data: <json>\n\n).
	StreamSSE StreamFormat = "sse"
	// StreamNDJSON is application/x-ndjson (one JSON object per line).
	// Same family as JSONL / NDJSON; not a single JSON array document.
	StreamNDJSON StreamFormat = "ndjson"
	// StreamJSONArray is application/json as one JSON array document:
	// [obj1,obj2,...] with items appended incrementally and closed by Close.
	StreamJSONArray StreamFormat = "json-array"
)

// Stream is a typed incremental response stream.
//
// Ideal handler usage:
//
//	stream, err := httpbind.NewStream[ChatEvent](w, r)
//	if err != nil { ... }
//	defer stream.Close()
//	_ = stream.Write(ChatEvent{Type: "delta", Delta: "hi"})
//	_ = stream.Write(ChatEvent{Type: "done"})
//
// Format (SSE vs NDJSON vs JSON array) is chosen once by rule:stream-content-negotiation.
// Write may be called many times; headers/status are sent only in NewStream.
// JSON array framing requires Close (via defer) so the trailing ']' is written.
type Stream[T any] struct {
	w       http.ResponseWriter
	format  StreamFormat
	enc     *json.Encoder
	closed  bool
	started bool // JSON array: '[' already written
}

// NewStream negotiates transport format from the request, writes response
// headers and 200 once, and returns a stream for incremental Write calls.
func NewStream[T any](w http.ResponseWriter, r *http.Request) (*Stream[T], error) {
	if w == nil {
		return nil, BadRequest(Problem{Code: "stream", Message: "nil ResponseWriter"})
	}
	format := NegotiateStreamFormat(r)
	s := &Stream[T]{
		w:      w,
		format: format,
		enc:    json.NewEncoder(w),
	}
	switch format {
	case StreamSSE:
		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		// Disable proxy buffering when supported (nginx etc.).
		w.Header().Set("X-Accel-Buffering", "no")
	case StreamJSONArray:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
	default: // NDJSON / JSONL
		w.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
	}
	w.WriteHeader(http.StatusOK)
	return s, nil
}

// Format returns the negotiated stream format (sse | ndjson | json-array).
func (s *Stream[T]) Format() StreamFormat {
	if s == nil {
		return ""
	}
	return s.format
}

// Write encodes one event in the negotiated format.
// Callable many times; does not re-send HTTP status or headers.
func (s *Stream[T]) Write(v T) error {
	if s == nil {
		return Internal(fmt.Errorf("httpbind: nil stream"))
	}
	if s.closed {
		return Internal(fmt.Errorf("httpbind: stream closed"))
	}
	switch s.format {
	case StreamSSE:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(s.w, "data: %s\n\n", data); err != nil {
			return err
		}
	case StreamJSONArray:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if !s.started {
			if _, err := s.w.Write([]byte{'['}); err != nil {
				return err
			}
			s.started = true
		} else {
			if _, err := s.w.Write([]byte{','}); err != nil {
				return err
			}
		}
		if _, err := s.w.Write(data); err != nil {
			return err
		}
	default: // NDJSON: Encoder.Encode appends '\n'
		if err := s.enc.Encode(v); err != nil {
			return err
		}
	}
	if f, ok := s.w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// Close marks the stream finished. Idempotent.
// For JSON array format, Close writes the trailing ']' (or "[]" if no Write).
// SSE and NDJSON do not require a special trailer; still call Close for symmetry.
func (s *Stream[T]) Close() error {
	if s == nil || s.closed {
		return nil
	}
	s.closed = true
	if s.format == StreamJSONArray {
		var err error
		if !s.started {
			_, err = s.w.Write([]byte("[]"))
		} else {
			_, err = s.w.Write([]byte{']'})
		}
		if err != nil {
			return err
		}
		if f, ok := s.w.(http.Flusher); ok {
			f.Flush()
		}
	}
	return nil
}

// NegotiateStreamFormat selects SSE, NDJSON, or JSON array using:
//  1. ?stream= query
//  2. Accept
//  3. User-Agent heuristics
//  4. default NDJSON
//
// Exported for tests and advanced callers.
//
// Note: NDJSON/JSONL (line-delimited objects) is distinct from JSON array
// (a single [...] document). application/json selects the array form;
// application/x-ndjson / application/jsonl select NDJSON.
func NegotiateStreamFormat(r *http.Request) StreamFormat {
	if r == nil {
		return StreamNDJSON
	}

	// 1) stream query parameter
	if q := strings.TrimSpace(r.URL.Query().Get("stream")); q != "" {
		switch strings.ToLower(q) {
		case "sse", "event-stream", "events", "eventstream":
			return StreamSSE
		case "ndjson", "jsonl", "nd", "lines":
			return StreamNDJSON
		case "json", "array", "json-array", "jsonarray":
			return StreamJSONArray
		}
	}

	// 2) Accept — first matching media type wins (left to right).
	if accept := r.Header.Get("Accept"); accept != "" {
		for part := range strings.SplitSeq(accept, ",") {
			media := strings.TrimSpace(strings.Split(part, ";")[0])
			media = strings.ToLower(media)
			switch media {
			case "text/event-stream":
				return StreamSSE
			case "application/x-ndjson", "application/ndjson", "application/jsonl":
				return StreamNDJSON
			case "application/json":
				// Full JSON array document (not JSONL).
				return StreamJSONArray
			}
		}
	}

	// 3) User-Agent
	ua := strings.ToLower(r.Header.Get("User-Agent"))
	if ua != "" {
		if isBrowserUA(ua) {
			return StreamSSE
		}
		if strings.Contains(ua, "curl") || strings.Contains(ua, "wget") || strings.Contains(ua, "httpie") {
			return StreamNDJSON
		}
	}

	// 4) default — curl-friendly NDJSON (JSONL-style lines)
	return StreamNDJSON
}

func isBrowserUA(ua string) bool {
	// Common browser tokens. Avoid matching "curl" which never appears here.
	return strings.Contains(ua, "mozilla/") ||
		strings.Contains(ua, "chrome/") ||
		strings.Contains(ua, "safari/") ||
		strings.Contains(ua, "firefox/") ||
		strings.Contains(ua, "edg/") ||
		strings.Contains(ua, "applewebkit")
}
