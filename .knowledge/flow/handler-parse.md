---
id: flow:handler-parse
type: flow
title: Handler Parse Flow
---
Static analysis flow from route registration through wrapper unwrap to Bind, Write, and error discovery.

```yaml
flow:
  trigger: same-package net/http route registration
  steps:
    - id: typecheck
      action: load package with go/types for rule:go-types-symbol-identity
    - id: discover-route
      action: match concept:route-discovery patterns only for allowed net/http symbols
    - id: unwrap-wrappers
      action: apply concept:stdlib-wrapper-unwrap and rule:nested-wrapper-unwrap
    - id: unwrap-middleware
      action: best-effort rule:custom-middleware-unwrap
    - id: resolve-handler
      action: resolve concept:handler-forms target
    - id: discover-request
      action: rule:request-model-discovery via api:bind (types-resolved)
    - id: discover-response
      action: rule:response-model-discovery via api:write / api:new-stream (types-resolved)
    - id: discover-errors
      action: rule:error-response-discovery via concept:error-helpers (types-resolved)
    - id: collect-route-metadata
      action: record wrapper metadata for concept:openapi-generation
  failure:
    unsupported_registration: rule:unsupported-route-patterns
symbol_identity: rule:go-types-symbol-identity
related:
  - concept:handler-discovery
  - flow:code-generation
  - requirement:strict-symbol-identity
  - rule:go-types-symbol-identity
```

