---
id: api:decode-json
type: api
title: jsonbind.DecodeJSON
---
Generic JSON decoder from io.Reader into typed T; independent of *http.Request.

```yaml
signature: "func DecodeJSON[T any](r io.Reader) (T, error)"
example: "v, err := jsonbind.DecodeJSON[CreateUserResponse](r)"
pair: api:encode-json
behavior:
  - decode one JSON value from r into T
  - compact JSON only; no pretty-print options
  - not HTTP: no Content-Type check, no query/path/header bind
  - on success return T; on failure return zero T and error
  - prefer generated codec when T is registered (decision:reflection-free)
  - enforce policy:json-read-limit before retaining the complete document
limit_api: "func DecodeJSONLimit[T any](r io.Reader, limit int64) (T, error)"
errors:
  - invalid JSON: transport-neutral jsonbind.Error
  - invalid field: jsonbind.Error with field identity and cause
  - missing codec for T: missing_codec jsonbind.Error
  - input over limit: payload_too_large jsonbind.Error
http_mapping: api:bind converts JSON errors to HTTP validation, 400, or 413 errors
differs_from:
  api:bind: Bind maps full HTTP request; DecodeJSON is document-only
  ReadJSONMap: internal map[string]json.RawMessage helper for generated binders
related:
  - concept:standalone-json-codec
  - api:encode-json
  - api:bind
  - concept:code-generation
  - system:tinybind
  - policy:json-read-limit
```
