package httpbinder

import (
	"fmt"
	"io"
)

// DecodeJSON decodes one JSON value from r into T using a generated codec.
// It does not inspect HTTP headers or use reflection on T's fields.
func DecodeJSON[T any](r io.Reader) (T, error) {
	var zero T
	fn, ok := lookupDecoder(typeKey[T]())
	if !ok {
		return zero, missingDecoderError(typeKey[T]())
	}
	if r == nil {
		return zero, BadRequest(Problem{Code: "json_parse", Message: "nil reader"}, nil)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return zero, BadRequest(Problem{Code: "body_read", Message: "failed to read JSON"}, err)
	}
	out, err := fn(data)
	if err != nil {
		return zero, err
	}
	return out.(T), nil
}

// EncodeJSON encodes v as compact JSON to w using a generated codec.
// It does not set HTTP headers or status.
func EncodeJSON[T any](w io.Writer, v T) error {
	fn, ok := lookupEncoder(typeKey[T]())
	if !ok {
		return missingEncoderError(typeKey[T]())
	}
	if w == nil {
		return Internal(fmt.Errorf("httpbinder: nil writer"))
	}
	return fn(w, v)
}
