---
id: requirement:product-goals
type: requirement
title: Product Goals
---
Core product goals for httpbinder API design and runtime behavior.

```yaml
goals:
  - Go-first API development
  - reflection-free runtime
  - code generation only
  - unified JSON / form / multipart handling
  - browser-friendly APIs
  - curl-friendly APIs
  - RFC 9457 compliant errors
  - automatic OpenAPI generation
  - TinyGo compatible
  - type-safe streaming APIs
maps_to:
  - decision:reflection-free
  - decision:single-source-of-truth
  - decision:stdlib-servemux
  - concept:request-binding
  - concept:streaming
  - concept:net-http-handler
  - policy:problem-details
  - concept:openapi-generation
  - requirement:tinygo-wasm
  - api:bind
  - api:write
  - api:write-error
```
