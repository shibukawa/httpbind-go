---
id: decision:single-source-of-truth
type: decision
title: Single Source of Truth
---
Developers define Go types only; all binders, writers, validation, OpenAPI, and streaming metadata are generated.

```yaml
authority: Go types
not_authority:
  - OpenAPI document as primary input
pipeline:
  - from: Go types
    to: flow:code-generation
  - from: flow:code-generation
    artifacts:
      - request binder
      - response writer
      - validation
      - OpenAPI
      - streaming metadata
runtime: generated code is the implementation
related:
  - vision:tinybind
  - concept:code-generation
  - concept:openapi-generation
  - concept:openapi-embed
  - requirement:openapi-goals
```
