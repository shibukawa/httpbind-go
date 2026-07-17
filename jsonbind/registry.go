package jsonbind

import (
	"fmt"
	"io"
	"reflect"
	"sync"
)

func typeKey[T any]() reflect.Type { return reflect.TypeFor[T]() }

type decodeFunc func([]byte) (any, error)
type encodeFunc func(io.Writer, any) error

var decoders sync.Map
var encoders sync.Map

// RegisterDecode registers a generated JSON document decoder for T.
func RegisterDecode[T any](fn func([]byte) (T, error)) {
	decoders.Store(typeKey[T](), decodeFunc(func(data []byte) (any, error) { return fn(data) }))
}

// RegisterEncode registers a generated compact JSON encoder for T.
func RegisterEncode[T any](fn func(io.Writer, T) error) {
	encoders.Store(typeKey[T](), encodeFunc(func(w io.Writer, v any) error { return fn(w, v.(T)) }))
}

func lookupDecoder(t reflect.Type) (decodeFunc, bool) {
	v, ok := decoders.Load(t)
	if !ok {
		return nil, false
	}
	return v.(decodeFunc), true
}

func lookupEncoder(t reflect.Type) (encodeFunc, bool) {
	v, ok := encoders.Load(t)
	if !ok {
		return nil, false
	}
	return v.(encodeFunc), true
}

func missingDecoderError(t reflect.Type) error {
	return newError("missing_codec", fmt.Sprintf("jsonbind: no JSON decoder registered for %s", t), nil)
}

func missingEncoderError(t reflect.Type) error {
	return newError("missing_codec", fmt.Sprintf("jsonbind: no JSON encoder registered for %s", t), nil)
}
