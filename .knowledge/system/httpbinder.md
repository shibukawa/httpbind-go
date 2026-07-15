---
id: system:httpbinder
type: system
title: httpbinder Library
---
Go library and code generator for typed HTTP request binding, response writing, validation, streaming, and OpenAPI.

```yaml
role: runtime plus ahead-of-time generator
runtime_style: generated code only; no reflection
public_api:
  - api:bind
  - api:write
  - api:write-error
  - concept:error-helpers
primary_inputs:
  - developer-defined Go types
  - struct field tags for source restriction
  - same-package net/http handlers
outputs:
  - generated bind and write functions
  - validation code
  - OpenAPI schemas
  - streaming metadata
related:
  - vision:httpbinder
  - flow:code-generation
  - flow:handler-request
  - concept:request-binding
  - concept:response-binding
  - concept:streaming
  - concept:net-http-handler
  - concept:handler-discovery
  - flow:handler-parse
  - concept:stdlib-wrapper-unwrap
  - policy:problem-details
  - decision:stdlib-servemux
  - concept:openapi-generation
  - concept:openapi-embed
  - api:openapi-json
  - api:openapi-yaml
```
