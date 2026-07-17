---
id: flow:handler-request
type: flow
title: Handler Request Flow
---
Typical request path through stdlib handler, bind, service, and write or write-error.

```yaml
flow:
  trigger: client HTTP request to registered ServeMux route
  steps:
    - id: route
      action: match decision:stdlib-servemux pattern
    - id: bind
      action: api:bind typed request
      on_error: api:write-error
    - id: service
      action: call concept:service-layer with context and request
      on_error: api:write-error
    - id: write
      action: api:write typed response or concept:streaming
  failure:
    default: api:write-error with policy:problem-details
related:
  - concept:net-http-handler
  - system:tinybind
```
