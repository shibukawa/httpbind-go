package httpbinder

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

// DefaultMultipartMaxMemory is the maxMemory argument passed to
// http.Request.ParseMultipartForm for generated binders (32 MiB).
// Larger file parts spill to temporary files; request body caps are separate
// (e.g. http.MaxBytesReader / MaxBytesHandler).
const DefaultMultipartMaxMemory int64 = 32 << 20

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

// IsMultipartRequest reports multipart/form-data.
func IsMultipartRequest(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	media, _, _ := strings.Cut(ct, ";")
	media = strings.TrimSpace(strings.ToLower(media))
	return media == "multipart/form-data"
}

// ParseMultipartMap parses a multipart/form-data body into scalar form fields
// (first value wins) and named file parts (first file wins per field name).
// Oversized bodies (MaxBytesReader / message-too-large) map to HTTP 413.
func ParseMultipartMap(r *http.Request) (form map[string]string, files map[string]File, err error) {
	if err := r.ParseMultipartForm(DefaultMultipartMaxMemory); err != nil {
		return nil, nil, multipartParseError(err)
	}
	form = make(map[string]string)
	files = make(map[string]File)
	if r.MultipartForm == nil {
		return form, files, nil
	}
	for k, vs := range r.MultipartForm.Value {
		if len(vs) > 0 {
			form[k] = vs[0]
		}
	}
	for k, fhs := range r.MultipartForm.File {
		if len(fhs) == 0 {
			continue
		}
		f, err := fileFromHeader(fhs[0])
		if err != nil {
			return nil, nil, BindError(k, "payload", "unreadable file")
		}
		files[k] = f
	}
	return form, files, nil
}

func fileFromHeader(fh *multipart.FileHeader) (File, error) {
	rc, err := fh.Open()
	if err != nil {
		return File{}, err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return File{}, err
	}
	ct := fh.Header.Get("Content-Type")
	size := fh.Size
	if size <= 0 {
		size = int64(len(data))
	}
	return File{
		Filename:    fh.Filename,
		ContentType: ct,
		Size:        size,
		Content:     data,
	}, nil
}

func multipartParseError(err error) error {
	if isRequestTooLarge(err) {
		return PayloadTooLarge(Problem{Code: "payload_too_large", Message: "multipart body too large"}, err)
	}
	return BadRequest(Problem{Code: "multipart_parse", Message: "invalid multipart body"}, err)
}

// isRequestTooLarge reports body/message size limit errors without errors.As,
// matching AsHTTPError's TinyGo-friendly unwrap style.
func isRequestTooLarge(err error) bool {
	for err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			return true
		}
		if err == multipart.ErrMessageTooLarge {
			return true
		}
		msg := err.Error()
		if strings.Contains(msg, "request body too large") ||
			strings.Contains(msg, "message too large") ||
			strings.Contains(msg, "http: request body too large") {
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
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
