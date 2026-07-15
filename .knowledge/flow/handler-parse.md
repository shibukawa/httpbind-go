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
    - id: discover-route
      action: match concept:route-discovery patterns
    - id: unwrap-wrappers
      action: apply concept:stdlib-wrapper-unwrap and rule:nested-wrapper-unwrap
    - id: unwrap-middleware
      action: best-effort rule:custom-middleware-unwrap
    - id: resolve-handler
      action: resolve concept:handler-forms target
    - id: discover-request
      action: rule:request-model-discovery via api:bind
    - id: discover-response
      action: rule:response-model-discovery via api:write
    - id: discover-errors
      action: rule:error-response-discovery via concept:error-helpers
    - id: collect-route-metadata
      action: record wrapper metadata for concept:openapi-generation
  failure:
    unsupported_registration: rule:unsupported-route-patterns
related:
  - concept:handler-discovery
  - flow:code-generation
```
