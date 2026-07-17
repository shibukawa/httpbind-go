---
id: rule:same-package-convention
type: rule
title: Same-Package Analysis Convention
---
Intended app convention: keep registration, handlers, and request/response models in one package; prefer literal patterns and mux-level middleware.

```yaml
status: implemented
requirement: requirement:analysis-diagnostics
conventions:
  route_registration_and_handler: same Go package
  request_and_response_structs: same Go package as handlers
  route_pattern: static string literal only
  middleware: prefer wrapping whole ServeMux (or shared chain) over per-route opaque wrappers
rationale:
  - no cross-package handler body analysis required
  - keeps go/types discovery tractable
  - OpenAPI stays aligned with source when conventions hold
out_of_scope_analysis:
  - dto.Request in another package
  - handlers defined only in another package
  - dynamic pattern concatenation or table-driven registration without literals
  - pointer type args and complex nested generic type args beyond supported forms
related:
  - requirement:analysis-diagnostics
  - rule:analysis-diagnostics-check
  - rule:unsupported-route-patterns
  - concept:handler-discovery
  - decision:stdlib-servemux
```
