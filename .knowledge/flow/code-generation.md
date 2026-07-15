---
id: flow:code-generation
type: flow
title: Code Generation Pipeline
---
Generator reads same-package handlers and Go types, then emits runtime bind/write functions, validation, OpenAPI, and streaming metadata from one IR.

```yaml
flow:
  trigger: developer defines Go types and net/http handlers
  steps:
    - id: discover-handlers
      action: run flow:handler-parse on same-package registrations
      refs:
        - concept:handler-discovery
        - concept:route-discovery
        - decision:stdlib-servemux
    - id: unwrap-wrappers
      action: unwrap stdlib wrappers and custom middleware when static
      refs:
        - concept:stdlib-wrapper-unwrap
        - rule:nested-wrapper-unwrap
        - rule:custom-middleware-unwrap
    - id: discover-models
      action: detect Bind/Write/error constructors in handler bodies
      refs:
        - rule:request-model-discovery
        - rule:response-model-discovery
        - rule:error-response-discovery
    - id: parse-go-types
      action: analyze discovered struct fields and tags
    - id: build-ir
      action: build shared intermediate representation including route metadata
    - id: emit-binders
      action: generate bind* functions for request types
      refs:
        - concept:request-binding
        - api:bind
    - id: emit-writers
      action: generate write* functions for response and stream types
      refs:
        - concept:response-binding
        - concept:streaming
        - api:write
    - id: emit-validation
      action: generate validation logic
    - id: emit-streaming-metadata
      action: generate streaming transport metadata
      refs:
        - concept:streaming
    - id: emit-openapi
      action: generate OpenAPI 3.1 model, embed, and serve handlers
      refs:
        - concept:openapi-generation
        - concept:openapi-embed
        - api:openapi-json
        - api:openapi-yaml
        - decision:openapi-31
  invariant: all artifacts derive from the same IR
  related:
    - decision:single-source-of-truth
    - system:httpbinder
    - concept:code-generation
    - flow:handler-parse
    - requirement:openapi-goals
```
