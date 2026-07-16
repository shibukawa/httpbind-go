---
id: system:httpbinder
type: system
title: httpbinder Library
---
Go library and code generator for typed HTTP request binding, response writing, validation, streaming, OpenAPI, and standalone JSON I/O.

```yaml
role: runtime plus ahead-of-time generator
runtime_style: generated code only; no reflection
public_api:
  - api:bind
  - api:write
  - api:write-error
  - api:decode-json
  - api:encode-json
  - concept:error-helpers
  - concept:standalone-json-codec
primary_inputs:
  - developer-defined Go types
  - struct field tags for source restriction
  - same-package net/http handlers
  - io.Reader / io.Writer for standalone JSON
outputs:
  - generated bind and write functions
  - validation code
  - OpenAPI schemas
  - streaming metadata
  - typed JSON codecs for registered models
related:
  - vision:httpbinder
  - flow:code-generation
  - flow:handler-request
  - concept:request-binding
  - concept:response-binding
  - concept:standalone-json-codec
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
  - api:decode-json
  - api:encode-json
```

