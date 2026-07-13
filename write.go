package httpbinder

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Write serializes a typed response value to the HTTP response via a registered writer.
func Write[T any](w http.ResponseWriter, r *http.Request, value T) error {
	fn, ok := lookupWriter(typeKey[T]())
	if !ok {
		return missingWriterError(typeKey[T]())
	}
	return fn(w, r, value)
}

// WriteError writes err as an RFC 9457 Problem Details response.
// Internal causes are not exposed in the client body.
//
// JSON is written without encoding/json for the problem document so TinyGo
// does not hit unimplemented reflect.AssignableTo when binders also use
// json.RawMessage (a known interaction in TinyGo's encoding/json).
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	_ = r
	status := http.StatusInternalServerError
	title := "Internal Server Error"
	detail := "internal error"
	code := "internal"
	var fields []FieldError

	if he, ok := AsHTTPError(err); ok {
		status = he.Status
		if he.Title != "" {
			title = he.Title
		} else {
			title = http.StatusText(status)
		}
		if he.Problem.Message != "" {
			detail = he.Problem.Message
		} else {
			detail = title
		}
		if he.Problem.Code != "" {
			code = he.Problem.Code
		}
		// Hide internal implementation details from clients for 5xx.
		if status >= 500 {
			detail = title
			code = "internal"
		}
		fields = he.Fields
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_, _ = w.Write(encodeProblemJSON(title, detail, code, status, fields))
}

func encodeProblemJSON(title, detail, code string, status int, fields []FieldError) []byte {
	var b strings.Builder
	b.WriteString(`{"type":"about:blank","title":`)
	b.WriteString(strconv.Quote(title))
	b.WriteString(`,"status":`)
	b.WriteString(strconv.Itoa(status))
	b.WriteString(`,"detail":`)
	b.WriteString(strconv.Quote(detail))
	b.WriteString(`,"code":`)
	b.WriteString(strconv.Quote(code))
	if len(fields) > 0 {
		b.WriteString(`,"errors":[`)
		for i, f := range fields {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(`{"field":`)
			b.WriteString(strconv.Quote(f.Field))
			b.WriteString(`,"location":`)
			b.WriteString(strconv.Quote(f.Location))
			b.WriteString(`,"message":`)
			b.WriteString(strconv.Quote(f.Message))
			b.WriteByte('}')
		}
		b.WriteByte(']')
	}
	b.WriteByte('}')
	return []byte(b.String())
}

// WriteJSON is a helper for generated writers: encode a pre-built map/slice without
// reflecting over application structs. Content-Type is application/json.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
