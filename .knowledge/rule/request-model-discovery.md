---
id: rule:request-model-discovery
type: rule
title: Request Model Discovery
---
Request models are discovered from the generic type argument of httpbinder.Bind[T](r).

```yaml
detection_call: "httpbinder.Bind[T](r)"
example: "input, err := httpbinder.Bind[CreateUserRequest](r)"
model_source: generic type argument T
related:
  - api:bind
  - concept:request-binding
  - concept:handler-discovery
  - concept:openapi-generation
```
