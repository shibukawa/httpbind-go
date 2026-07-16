---
id: api:encode-json
type: api
title: httpbinder.EncodeJSON
---
Generic compact JSON encoder of typed T to io.Writer; independent of http.ResponseWriter.

```yaml
status: implemented
signature: "func EncodeJSON[T any](w io.Writer, v T) error"
example: "err := httpbinder.EncodeJSON(w, output)"
pair: api:decode-json
behavior:
  - encode v as compact JSON to w
  - no pretty-print / indent API
  - not HTTP: no Status Code, no Content-Type header
  - prefer generated codec when T is registered (decision:reflection-free)
name_note: |
  Existing untyped helper WriteJSON(http.ResponseWriter, status int, v any) stays
  an HTTP/codegen helper name; public standalone API is EncodeJSON (no rename required).
differs_from:
  api:write: Write sets HTTP status/headers and uses Accept/stream paths
  api:write-error: problem+json for errors
  WriteJSON: internal/generated HTTP helper writing application/json responses
related:
  - concept:standalone-json-codec
  - api:decode-json
  - api:write
  - concept:code-generation
  - concept:response-binding
  - system:httpbinder
```
