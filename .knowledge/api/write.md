---
id: api:write
type: api
title: httpbinder.Write
---
Generic response writer that serializes a typed value or stream to the HTTP response.

```yaml
signature: "func Write[T any](w http.ResponseWriter, r *http.Request, value T) error"
examples:
  - "httpbinder.Write[CreateUserResponse](w, r, output)"
behavior:
  - serialize ordinary response values with HTTP 200 OK
  - no runtime field reflection on T
  - streaming uses api:new-stream (not Write[Stream[T]] for incremental handlers)
status: always 200 unless using api:write-status
uses:
  - concept:response-binding
  - concept:code-generation
discovery: rule:response-model-discovery
openapi: rule:openapi-success-response
related:
  - system:httpbinder
  - concept:net-http-handler
  - concept:handler-discovery
  - api:write-error
  - api:write-status
  - rule:openapi-success-status
```

