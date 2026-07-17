---
id: requirement:analysis-diagnostics
type: requirement
title: Diagnostics for Undiscoverable Analysis
---
Undiscoverable routes and models must not be dropped silently from OpenAPI; host check surfaces error or warning.

```yaml
status: implemented
problem:
  - dynamic route patterns ignored without report
  - cross-package handlers ignored without report
  - unanalyzable middleware may hide the leaf handler without report
  - cross-package request/response types ignored without report
  - pointer or complex generic type args ignored without report
risk: silent OpenAPI drift from real HTTP surface
policy: rule:same-package-convention
diagnostics: rule:analysis-diagnostics-check
not_required:
  - full cross-package handler analysis
  - deep DI / dynamic router support
related:
  - rule:same-package-convention
  - rule:analysis-diagnostics-check
  - rule:unsupported-route-patterns
  - concept:handler-discovery
  - concept:route-discovery
  - concept:openapi-generation
  - flow:handler-parse
```
