---
id: rule:wrapper-max-bytes-handler
type: rule
title: MaxBytesHandler Unwrap
---
Unwrap http.MaxBytesHandler, analyze inner handler, and record max request body bytes.

```yaml
wrapper: "http.MaxBytesHandler(h, n)"
example: |
  mux.Handle(
      "POST /upload",
      http.MaxBytesHandler(
          http.HandlerFunc(uploadHandler),
          10<<20,
      ),
  )
parse:
  - unwrap h
  - analyze inner handler
  - record max_request_body_bytes = n
openapi:
  - document maximum request body size metadata
  - may add 413 Payload Too Large as possible error response
related:
  - concept:stdlib-wrapper-unwrap
  - rule:standard-error-mapping
  - concept:openapi-generation
```
