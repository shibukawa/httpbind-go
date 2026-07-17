package httpbind

import (
	"fmt"
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

var (
	binders sync.Map // reflect.Type -> binderFunc
	writers sync.Map // reflect.Type -> writerFunc
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

func missingBinderError(t reflect.Type) error {
	return Internal(fmt.Errorf("httpbind: no binder registered for %s", t.String()))
}

func missingWriterError(t reflect.Type) error {
	return Internal(fmt.Errorf("httpbind: no writer registered for %s", t.String()))
}
