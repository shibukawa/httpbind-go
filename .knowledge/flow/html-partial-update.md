---
id: flow:html-partial-update
type: flow
title: HTML Partial Update Flow
---
Client navigation flow that omits unchanged component boundary HTML.

```yaml
flow:
  trigger: browser changes path or search parameters without full navigation
  steps:
    - id: request
      action: send target URL and prior data:component-update-manifest validators
    - id: execute
      action: rerun data:html-render-route-plan, requirement:nested-layout-composition chain, and page with current request context
    - id: retain-layouts
      action: preserve requirement:layout-reuse-boundaries whose frame validators remain unchanged
    - id: identify
      action: build instances through rule:component-instance-identity and requirement:partial-update-boundaries
    - id: reuse
      action: consult requirement:component-output-cache for eligible instance input validators
    - id: compare
      action: compare current content validators with client validators
    - id: respond
      output: data:component-delta-response containing changed and structural operations only
    - id: apply
      action: fixed browser runtime applies operations and stores next data:component-update-manifest
  failure:
    incompatible_or_ambiguous: fall back to complete HTML navigation
```
