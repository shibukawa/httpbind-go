package httpbinder

import "net/http"

// Bind maps an HTTP request into a typed request value.
// Dispatch uses a registry of generated binders; field mapping does not use reflect.
func Bind[T any](r *http.Request) (T, error) {
	var zero T
	fn, ok := lookupBinder(typeKey[T]())
	if !ok {
		return zero, missingBinderError(typeKey[T]())
	}
	out, err := fn(r)
	if err != nil {
		return zero, err
	}
	return out.(T), nil
}
