---
id: rule:wrapper-timeout-handler
type: rule
title: TimeoutHandler Unwrap
---
Unwrap http.TimeoutHandler, analyze inner handler, and record timeout duration and message when static.

```yaml
wrapper: "http.TimeoutHandler(h, dt, msg)"
example: |
  mux.Handle(
      "POST /jobs",
      http.TimeoutHandler(
          http.HandlerFunc(jobHandler),
          30*time.Second,
          "timeout",
      ),
  )
parse:
  - unwrap h
  - analyze inner handler
  - record timeout duration if statically known
  - record timeout message if statically known
openapi:
  - may add 503 Service Unavailable as possible error response
  - may document timeout behavior
related:
  - concept:stdlib-wrapper-unwrap
  - concept:openapi-generation
```
