---
id: rule:nested-wrapper-unwrap
type: rule
title: Nested Wrapper Unwrap
---
Nested stdlib wrappers are unwrapped outer-to-inner; metadata from each layer is collected.

```yaml
example: |
  mux.Handle(
      "POST /upload",
      http.TimeoutHandler(
          http.MaxBytesHandler(
              http.HandlerFunc(uploadHandler),
              10<<20,
          ),
          30*time.Second,
          "timeout",
      ),
  )
parse_order:
  - unwrap TimeoutHandler
  - unwrap MaxBytesHandler
  - unwrap HandlerFunc
  - analyze uploadHandler
collected_metadata_example:
  timeout: 30s
  max_request_body_bytes: 10485760
related:
  - concept:stdlib-wrapper-unwrap
  - rule:wrapper-timeout-handler
  - rule:wrapper-max-bytes-handler
```
