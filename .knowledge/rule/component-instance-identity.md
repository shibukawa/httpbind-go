---
id: rule:component-instance-identity
type: rule
title: Stable Component Instance Identity
---
Identify the same updateable component invocation across two executions without trusting browser-provided arguments.

```yaml
source: concept:html-render-runtime-extensions
identity_parts:
  - route and runtime layout-chain identity
  - generated component call-site identity
  - explicit stable item key when one call site produces repeated instances
properties:
  - deterministic for the same logical instance across search parameter changes
  - unique within one rendered document
  - independent from generated transport marker IDs
  - opaque and safe in HTML attributes and protocol fields
validation:
  - repeated update boundaries require a statically checked stable key expression
  - duplicate instance IDs are rendering errors before delta publication
  - browser IDs and validators are comparison hints, never authority for component arguments or access control
structural_change:
  missing_old_id: insertion or nearest-ancestor replacement
  missing_new_id: removal or nearest-ancestor replacement
  changed_parent_or_order: move operation or nearest-ancestor replacement
```
