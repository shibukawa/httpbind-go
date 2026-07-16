---
id: api:decode-json
type: api
title: httpbinder.DecodeJSON
---
Generic JSON decoder from io.Reader into typed T; independent of *http.Request.

```yaml
status: implemented
signature: "func DecodeJSON[T any](r io.Reader) (T, error)"
example: "v, err := httpbinder.DecodeJSON[CreateUserResponse](r)"
pair: api:encode-json
behavior:
  - decode one JSON value from r into T
  - compact JSON only; no pretty-print options
  - not HTTP: no Content-Type check, no query/path/header bind
  - on success return T; on failure return zero T and error
  - prefer generated codec when T is registered (decision:reflection-free)
errors:
  - invalid JSON: 400-style problem-capable error or plain decode error (implementation chooses HTTPError when useful)
  - missing codec for T if registry-only mode: clear internal/missing-codec error
differs_from:
  api:bind: Bind maps full HTTP request; DecodeJSON is document-only
  ReadJSONMap: internal map[string]json.RawMessage helper for generated binders
related:
  - concept:standalone-json-codec
  - api:encode-json
  - api:bind
  - concept:code-generation
  - system:httpbinder
```
