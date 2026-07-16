package httpbinder

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
)

// typeKey is used only for registry dispatch identity — not for field walking.
// Generated binders/writers perform all field mapping without reflect.
func typeKey[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

type binderFunc func(*http.Request) (any, error)
type writerFunc func(http.ResponseWriter, *http.Request, any) error
type decodeFunc func([]byte) (any, error)
type encodeFunc func(io.Writer, any) error

var (
	binders  sync.Map // reflect.Type -> binderFunc
	writers  sync.Map // reflect.Type -> writerFunc
	decoders sync.Map // reflect.Type -> decodeFunc
	encoders sync.Map // reflect.Type -> encodeFunc
)

// RegisterBind registers a generated binder for T.
// Call from generated init(); field mapping lives entirely inside fn.
func RegisterBind[T any](fn func(*http.Request) (T, error)) {
	binders.Store(typeKey[T](), binderFunc(func(r *http.Request) (any, error) {
		return fn(r)
	}))
}

// RegisterWrite registers a generated writer for T.
func RegisterWrite[T any](fn func(http.ResponseWriter, *http.Request, T) error) {
	writers.Store(typeKey[T](), writerFunc(func(w http.ResponseWriter, r *http.Request, v any) error {
		return fn(w, r, v.(T))
	}))
}

// RegisterDecode registers a generated JSON decoder for T (document only).
func RegisterDecode[T any](fn func([]byte) (T, error)) {
	decoders.Store(typeKey[T](), decodeFunc(func(data []byte) (any, error) {
		return fn(data)
	}))
}

// RegisterEncode registers a generated compact JSON encoder for T.
func RegisterEncode[T any](fn func(io.Writer, T) error) {
	encoders.Store(typeKey[T](), encodeFunc(func(w io.Writer, v any) error {
		return fn(w, v.(T))
	}))
}

func lookupBinder(t reflect.Type) (binderFunc, bool) {
	v, ok := binders.Load(t)
	if !ok {
		return nil, false
	}
	return v.(binderFunc), true
}

func lookupWriter(t reflect.Type) (writerFunc, bool) {
	v, ok := writers.Load(t)
	if !ok {
		return nil, false
	}
	return v.(writerFunc), true
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

func missingBinderError(t reflect.Type) error {
	return Internal(fmt.Errorf("httpbinder: no binder registered for %s", t.String()))
}

func missingWriterError(t reflect.Type) error {
	return Internal(fmt.Errorf("httpbinder: no writer registered for %s", t.String()))
}

func missingDecoderError(t reflect.Type) error {
	return Internal(fmt.Errorf("httpbinder: no JSON decoder registered for %s", t.String()))
}

func missingEncoderError(t reflect.Type) error {
	return Internal(fmt.Errorf("httpbinder: no JSON encoder registered for %s", t.String()))
}
