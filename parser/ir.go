package parser

import "encoding/json"

// Result is the structured parse output for a single Go package.
type Result struct {
	Routes []Route `json:"routes"`
}

// Route is one statically discovered net/http route registration.
type Route struct {
	Method   string        `json:"method"`
	Path     string        `json:"path"`
	Handler  Handler       `json:"handler"`
	Request  string        `json:"request,omitempty"`
	Response string        `json:"response,omitempty"`
	Stream   string        `json:"stream,omitempty"` // element type when response is Stream[T]
	Errors   []string      `json:"errors,omitempty"`
	Wrappers WrapperMeta   `json:"wrappers"`
}

// Handler describes how the leaf handler was expressed in source.
type Handler struct {
	Form string `json:"form"`           // named | inline | struct
	Name string `json:"name,omitempty"` // function or type name when known
}

// WrapperMeta holds statically known stdlib wrapper metadata.
type WrapperMeta struct {
	AllowQuerySemicolons bool   `json:"allow_query_semicolons"`
	MaxRequestBodyBytes  *int64 `json:"max_request_body_bytes,omitempty"`
	StrippedPrefix       string `json:"stripped_prefix,omitempty"`
	Timeout              string `json:"timeout,omitempty"`
	TimeoutMessage       string `json:"timeout_message,omitempty"`
}

// Normalize sorts routes and nested slices for stable golden comparisons.
func (r *Result) Normalize() {
	if r == nil {
		return
	}
	if r.Routes == nil {
		r.Routes = []Route{}
	}
	for i := range r.Routes {
		sortStrings(r.Routes[i].Errors)
	}
	sortRoutes(r.Routes)
}

// JSON returns canonical indented JSON for golden files.
func (r Result) JSON() ([]byte, error) {
	r.Normalize()
	return json.MarshalIndent(r, "", "  ")
}
