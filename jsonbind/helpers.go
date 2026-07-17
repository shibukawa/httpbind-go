package jsonbind

import (
	"encoding/json"
	"errors"
	"io"
	"sync/atomic"
)

// DefaultMaxJSONBodyBytes is the default JSON document limit (1 MiB).
const DefaultMaxJSONBodyBytes int64 = 1 << 20

var maxJSONBodyBytes atomic.Int64

// ErrBodyTooLarge reports that a JSON document exceeded its configured limit.
var ErrBodyTooLarge = errors.New("jsonbind: JSON body too large")

// SetMaxJSONBodyBytes changes the process-wide JSON document limit.
func SetMaxJSONBodyBytes(n int64) {
	if n <= 0 {
		maxJSONBodyBytes.Store(0)
		return
	}
	maxJSONBodyBytes.Store(n)
}

// MaxJSONBodyBytes returns the effective JSON document limit.
func MaxJSONBodyBytes() int64 {
	if n := maxJSONBodyBytes.Load(); n > 0 {
		return n
	}
	return DefaultMaxJSONBodyBytes
}

// ReadLimit reads at most limit bytes from r.
func ReadLimit(r io.Reader, limit int64) ([]byte, error) {
	if limit <= 0 {
		limit = DefaultMaxJSONBodyBytes
	}
	data, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, ErrBodyTooLarge
	}
	return data, nil
}

func readJSONBytes(r io.Reader, limit int64) ([]byte, error) { return ReadLimit(r, limit) }

// RawJSONMap decodes a JSON object into raw fields.
func RawJSONMap(raw json.RawMessage) (map[string]json.RawMessage, error) {
	if len(raw) == 0 {
		return map[string]json.RawMessage{}, nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, newError("json_parse", "invalid JSON object", err)
	}
	if m == nil {
		return nil, newError("json_parse", "JSON value must be an object", nil)
	}
	return m, nil
}

// BytesJSONMap decodes a complete JSON object document.
func BytesJSONMap(data []byte) (map[string]json.RawMessage, error) {
	return RawJSONMap(json.RawMessage(data))
}

// RawJSONArray decodes a JSON array into raw elements.
func RawJSONArray(raw json.RawMessage) ([]json.RawMessage, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil {
		return nil, newError("json_parse", "invalid JSON array", err)
	}
	return arr, nil
}

func decode[T any](raw json.RawMessage, message string) (T, error) {
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, newError("json_parse", message, err)
	}
	return out, nil
}

func DecodeJSONMapStringString(raw json.RawMessage) (map[string]string, error) {
	if len(raw) == 0 {
		return map[string]string{}, nil
	}
	m, err := decode[map[string]string](raw, "invalid string map")
	if m == nil && err == nil {
		m = map[string]string{}
	}
	return m, err
}
func DecodeJSONStringSlice(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	return decode[[]string](raw, "invalid string array")
}
func DecodeJSONIntSlice(raw json.RawMessage) ([]int, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	return decode[[]int](raw, "invalid int array")
}
func DecodeJSONInt64Slice(raw json.RawMessage) ([]int64, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	return decode[[]int64](raw, "invalid int64 array")
}
func DecodeJSONBoolSlice(raw json.RawMessage) ([]bool, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	return decode[[]bool](raw, "invalid bool array")
}
func DecodeJSONFloat64Slice(raw json.RawMessage) ([]float64, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	return decode[[]float64](raw, "invalid float64 array")
}
func DecodeJSONString(raw json.RawMessage) (string, error) {
	return decode[string](raw, "invalid string")
}
func DecodeJSONInt(raw json.RawMessage) (int, error) { return decode[int](raw, "invalid int") }
func DecodeJSONInt64(raw json.RawMessage) (int64, error) {
	return decode[int64](raw, "invalid int64")
}
func DecodeJSONBool(raw json.RawMessage) (bool, error) {
	return decode[bool](raw, "invalid bool")
}
func DecodeJSONFloat64(raw json.RawMessage) (float64, error) {
	return decode[float64](raw, "invalid float64")
}

// RestJSONAny returns JSON fields not named in exclude.
func RestJSONAny(body map[string]json.RawMessage, exclude []string) (map[string]any, error) {
	out := make(map[string]any)
	skip := excludeSet(exclude)
	for k, raw := range body {
		if skip[k] {
			continue
		}
		if len(raw) == 0 || string(raw) == "null" {
			out[k] = nil
			continue
		}
		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, newError("json_parse", "invalid JSON rest value", err)
		}
		out[k] = value
	}
	return out, nil
}

// RestJSONRaw copies JSON fields not named in exclude.
func RestJSONRaw(body map[string]json.RawMessage, exclude []string) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage)
	skip := excludeSet(exclude)
	for k, raw := range body {
		if skip[k] {
			continue
		}
		out[k] = append(json.RawMessage(nil), raw...)
	}
	return out
}

func excludeSet(exclude []string) map[string]bool {
	skip := make(map[string]bool, len(exclude))
	for _, key := range exclude {
		if key != "" && key != "*" {
			skip[key] = true
		}
	}
	return skip
}
