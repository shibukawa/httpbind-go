// Package jsonbind provides generated, reflection-free JSON document codecs.
package jsonbind

import "io"

// DecodeJSON decodes one JSON value from r into T using a generated codec.
// It does not inspect HTTP headers or use reflection on T's fields.
func DecodeJSON[T any](r io.Reader) (T, error) {
	var zero T
	fn, ok := lookupDecoder(typeKey[T]())
	if !ok {
		return zero, missingDecoderError(typeKey[T]())
	}
	if r == nil {
		return zero, newError("json_parse", "nil reader", nil)
	}
	data, err := readJSONBytes(r, MaxJSONBodyBytes())
	if err != nil {
		if err == ErrBodyTooLarge {
			return zero, newError("payload_too_large", "JSON body too large", err)
		}
		return zero, newError("body_read", "failed to read JSON", err)
	}
	out, err := fn(data)
	if err != nil {
		return zero, err
	}
	return out.(T), nil
}

// DecodeJSONLimit is DecodeJSON with a per-call byte limit. A non-positive
// limit uses MaxJSONBodyBytes.
func DecodeJSONLimit[T any](r io.Reader, limit int64) (T, error) {
	if limit <= 0 {
		limit = MaxJSONBodyBytes()
	}
	var zero T
	fn, ok := lookupDecoder(typeKey[T]())
	if !ok {
		return zero, missingDecoderError(typeKey[T]())
	}
	if r == nil {
		return zero, newError("json_parse", "nil reader", nil)
	}
	data, err := readJSONBytes(r, limit)
	if err == ErrBodyTooLarge {
		return zero, newError("payload_too_large", "JSON body too large", err)
	}
	if err != nil {
		return zero, newError("body_read", "failed to read JSON", err)
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
		return newError("internal", "jsonbind: nil writer", nil)
	}
	return fn(w, v)
}
