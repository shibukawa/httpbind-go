---
id: concept:stdlib-wrapper-unwrap
type: concept
title: Stdlib Handler Wrapper Unwrap
---
Generator unwraps selected net/http wrappers, continues to the inner handler, and records wrapper metadata for OpenAPI.

```yaml
supported_wrappers:
  - rule:wrapper-allow-query-semicolons
  - rule:wrapper-max-bytes-handler
  - rule:wrapper-strip-prefix
  - rule:wrapper-timeout-handler
nesting: rule:nested-wrapper-unwrap
custom_middleware: rule:custom-middleware-unwrap
behavior:
  - unwrap known wrapper
  - analyze inner handler
  - collect route metadata when statically known
related:
  - concept:handler-discovery
  - concept:openapi-generation
```
