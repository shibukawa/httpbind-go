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
  - serialize ordinary response values
  - no runtime reflection
  - streaming uses api:new-stream (not Write[Stream[T]] for incremental handlers)
uses:
  - concept:response-binding
  - concept:code-generation
discovery: rule:response-model-discovery
related:
  - system:httpbinder
  - concept:net-http-handler
  - concept:handler-discovery
  - api:write-error
```
