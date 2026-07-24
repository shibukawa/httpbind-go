---
id: decision:html-route-go-package-model
type: decision
title: HTML Route Go Package Model
---
Treat route subdirectories as template resources and emit the whole tree into one configured Go package.

```yaml
source:
  - concept:filesystem-html-routing
  - user package discussion 2026-07-23
review_gate: proposed package model requires user approval
route_tree:
  role: generator input hierarchy and URL namespace, not Go package hierarchy
  package_declaration: omitted in route templates or required to match the single configured output package
output:
  directory: one configured Go package outside or at the route resource root
  files: generator may split output mechanically by route while all files share one package
  package_name: configured once or inferred from the output directory
user_logic:
  location: any ordinary Go package chosen by the application; not required beside each route template
  binding: one injected data:html-route-dependencies value supplied to api:register-generated-html-routes
  imports: dependency implementation may import generated route types; generated routes never import the application implementation package
external_names:
  template_scope: local source name remains concise
  generated_scope: stable route-relative module identity prevents collisions across folders
  shared_logic: application may delegate multiple generated dependency groups to one service value
compatibility:
  - existing flat template mode retains same-directory package and external function behavior
  - route mode does not create implicit Go subpackages or require package declarations per folder
constraints:
  - no runtime reflection or string-based dependency lookup
  - generated dependency binding is complete and type-checked at Go compile time
  - avoid application import cycles by keeping implementation outside the generated route package
open_questions:
  - configuration syntax for resource root, output directory, and package name
  - whether route templates may include an optional shared package declaration
  - generated dependency interface versus function-field provider shape
```
