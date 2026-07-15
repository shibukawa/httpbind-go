---
id: requirement:tinygo-wasm
type: requirement
title: TinyGo and WASM Support
---
TinyGo is a first-class target; reflection-free generated code must work on TinyGo, WebAssembly, and embedded environments.

```yaml
priority: first-class
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
