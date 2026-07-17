---
id: requirement:tinygo-wasm
type: requirement
title: TinyGo and WASM Support
---
TinyGo is a first-class target; reflection-free generated code must work on TinyGo, WebAssembly, and embedded environments.

```yaml
priority: first-class
baseline:
  tinygo: 0.41.1
  go: 1.26.x
targets:
  - TinyGo
  - WebAssembly
  - embedded environments
depends_on:
  - decision:reflection-free
related:
  - vision:httpbinder
  - system:httpbinder
```
