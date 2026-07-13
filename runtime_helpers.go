package httpbinder

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Content-type helpers and scalar parsers used by generated binders.
// These do not inspect application struct fields via reflect.

// IsJSONRequest reports whether the request body should be treated as JSON.
func IsJSONRequest(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		return false
	}
	media, _, _ := strings.Cut(ct, ";")
	media = strings.TrimSpace(strings.ToLower(media))
	return media == "application/json"
}

// IsFormRequest reports application/x-www-form-urlencoded.
func IsFormRequest(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	media, _, _ := strings.Cut(ct, ";")
	media = strings.TrimSpace(strings.ToLower(media))
	return media == "application/x-www-form-urlencoded"
}

// ReadJSONMap decodes a JSON object body into a map of raw messages.
// Used by generated binders so they can pick named fields without reflect on T.
func ReadJSONMap(r *http.Request) (map[string]json.RawMessage, error) {
	if r.Body == nil {
		return map[string]json.RawMessage{}, nil
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, BadRequest(Problem{Code: "body_read", Message: "failed to read body"}, err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return map[string]json.RawMessage{}, nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, BadRequest(Problem{Code: "json_parse", Message: "invalid JSON body"}, err)
	}
	if m == nil {
		m = map[string]json.RawMessage{}
	}
	return m, nil
}

// ParseFormMap parses urlencoded form body into a flat map (first value wins).
func ParseFormMap(r *http.Request) (map[string]string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, BadRequest(Problem{Code: "form_parse", Message: "invalid form body"}, err)
	}
	out := make(map[string]string, len(r.PostForm))
	for k, vs := range r.PostForm {
		if len(vs) > 0 {
			out[k] = vs[0]
		}
	}
	return out, nil
}

// QueryValue returns the first query parameter value for key.
func QueryValue(r *http.Request, key string) (string, bool) {
	if r.URL == nil {
		return "", false
	}
	vs := r.URL.Query()[key]
	if len(vs) == 0 {
		return "", false
	}
	return vs[0], true
}

// PathValue returns the path value for key (Go 1.22+ ServeMux).
func PathValue(r *http.Request, key string) string {
	return r.PathValue(key)
}

// HeaderValue returns a request header.
func HeaderValue(r *http.Request, key string) string {
	return r.Header.Get(key)
}

// CookieValue returns a cookie value if present.
func CookieValue(r *http.Request, name string) (string, bool) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", false
	}
	return c.Value, true
}

// DecodeJSONString unmarshals a JSON raw value as string.
func DecodeJSONString(raw json.RawMessage) (string, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", err
	}
	return s, nil
}

// DecodeJSONInt unmarshals a JSON raw value as int.
func DecodeJSONInt(raw json.RawMessage) (int, error) {
	var n int
	if err := json.Unmarshal(raw, &n); err != nil {
		return 0, err
	}
	return n, nil
}

// DecodeJSONInt64 unmarshals a JSON raw value as int64.
func DecodeJSONInt64(raw json.RawMessage) (int64, error) {
	var n int64
	if err := json.Unmarshal(raw, &n); err != nil {
		return 0, err
	}
	return n, nil
}

// DecodeJSONBool unmarshals a JSON raw value as bool.
func DecodeJSONBool(raw json.RawMessage) (bool, error) {
	var b bool
	if err := json.Unmarshal(raw, &b); err != nil {
		return false, err
	}
	return b, nil
}

// DecodeJSONFloat64 unmarshals a JSON raw value as float64.
func DecodeJSONFloat64(raw json.RawMessage) (float64, error) {
	var f float64
	if err := json.Unmarshal(raw, &f); err != nil {
		return 0, err
	}
	return f, nil
}

// ParseInt converts a string to int.
func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ParseInt64 converts a string to int64.
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseBool converts a string to bool.
func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ParseFloat64 converts a string to float64.
func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
