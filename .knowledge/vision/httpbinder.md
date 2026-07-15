---
id: vision:httpbinder
type: vision
title: httpbinder Vision
---
httpbinder is a code-generation-first library that bridges Go types and HTTP APIs without runtime reflection.

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
  - system:httpbinder
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
