---
id: vision:tinybind
type: vision
title: tinybind-go Vision
---
tinybind-go provides code-generated typed binding without application-field reflection and isolates platform dependencies by mapping mode.

```yaml
source_of_truth:
  - Go types only
generated_from_types:
  - request binding
  - response serialization
  - streaming responses
  - error handling
  - validation via concept:check-validation
  - OpenAPI generation
principles:
  - decision:single-source-of-truth
  - decision:reflection-free
targets:
  - system:tinybind
  - requirement:tinygo-wasm
  - concept:code-generation
  - concept:openapi-generation
  - concept:net-http-handler
  - decision:stdlib-servemux
public_runtime:
  - api:bind
  - api:write
  - api:write-error
```
