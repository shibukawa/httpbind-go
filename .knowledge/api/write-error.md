---
id: api:write-error
type: api
title: httpbinder.WriteError
---
Writes an error as an RFC 9457 Problem Details response and keeps internal causes out of the client body.

```yaml
signature: "func WriteError(w http.ResponseWriter, r *http.Request, err error)"
example: |
  if err != nil {
      httpbinder.WriteError(w, r, err)
      return
  }
behavior:
  - resolve HTTP status
  - convert error to RFC 9457 response
  - log wrapped internal cause
  - hide internal implementation details from clients
policy: policy:problem-details
helpers: concept:error-helpers
related:
  - system:httpbinder
  - concept:net-http-handler
  - rule:error-cause-preservation
```
