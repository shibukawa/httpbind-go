---
id: concept:handler-discovery
type: concept
title: Handler Discovery
---
Generator analyzes same-package net/http route registration, handlers, selected stdlib wrappers, and Bind/Write/error call sites.

```yaml
scope: same Go package only
supports:
  - standard net/http handlers
  - http.HandlerFunc
  - handler structs with ServeHTTP
  - selected built-in net/http wrappers
out_of_scope:
  - cross-package handler implementation analysis
convention: rule:same-package-convention
diagnostics: rule:analysis-diagnostics-check
requirement_diagnostics: requirement:analysis-diagnostics
symbol_identity: rule:go-types-symbol-identity
requirement: requirement:strict-symbol-identity
pipeline:
  - concept:route-discovery
  - concept:handler-forms
  - concept:stdlib-wrapper-unwrap
  - rule:custom-middleware-unwrap
  - rule:request-model-discovery
  - rule:response-model-discovery
  - rule:error-response-discovery
unsupported: rule:unsupported-route-patterns
related:
  - flow:code-generation
  - concept:code-generation
  - concept:net-http-handler
  - decision:stdlib-servemux
  - concept:openapi-generation
  - rule:go-types-symbol-identity
  - requirement:strict-symbol-identity
  - rule:same-package-convention
  - rule:analysis-diagnostics-check
  - requirement:analysis-diagnostics
```


