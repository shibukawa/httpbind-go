---
id: decision:reflection-free
type: decision
title: Reflection-Free Runtime
---
Reflection is intentionally unsupported; binding and serialization use generated code only.

```yaml
status: accepted
forbidden:
  - runtime reflection
  - runtime tag parsing
rationale:
  - better performance
  - smaller binaries
  - TinyGo compatibility
  - WASM compatibility
  - compile-time validation
  - OpenAPI consistency with runtime
related:
  - vision:tinybind
  - requirement:tinygo-wasm
  - concept:code-generation
  - decision:single-source-of-truth
```
