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
verified:
  host_tinygo: runtime and generated mapping checks pass
  js_wasm_json: JSON-only generated fixture builds with tinygo build -target wasm
  js_wasm_http: root HTTP runtime fails while compiling net/http roundtrip_js.go
v0_1_3_boundary:
  solved: usage-directed output omits net/http for DecodeJSON / EncodeJSON-only generation
  remaining: importing the root runtime still compiles registry.go and other net/http files
resolution:
  module: github.com/shibukawa/tinybind-go
  JSON-only runtime: jsonbind
  HTTP runtime: httpbind root package
  SQL runtime: sqlbind
  dependency_check: jsonbind dependency graph excludes net/http and database/sql
acceptance:
  - importing JSON-only runtime code does not compile net/http or database/sql
  - JSON-only generated code imports only the JSON runtime boundary
  - host Go behavior and deterministic generation remain unchanged
targets:
  - TinyGo
  - WebAssembly
  - embedded environments
depends_on:
  - decision:reflection-free
  - decision:runtime-package-boundaries
related:
  - vision:tinybind
  - system:tinybind
```
