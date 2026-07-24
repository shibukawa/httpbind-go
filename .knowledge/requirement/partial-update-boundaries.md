---
id: requirement:partial-update-boundaries
type: requirement
title: Explicit Partial Update Boundaries
---
Mark components whose rendered DOM and validators may participate in client-side partial updates.

```yaml
source: concept:html-render-runtime-extensions
declaration:
  activation:
    explicit: update flag on an ordinary component
    automatic: requirement:layout-reuse-boundaries for generated route layouts
  cache_relation: independent from requirement:component-output-cache
rendering:
  - emit stable boundary markers using rule:component-instance-identity
  - add one entry to data:component-update-manifest
  - preserve normal complete HTML for initial navigation and non-update clients
parameters:
  source: generated page execution from current path, search parameters, request state, and typed parent inputs
  client_state: raw arguments need not be exposed or resubmitted
direct_update:
  capability: requirement:boundary-parameter-updates
  mutable_parameters: rule:client-mutable-component-parameters
  context: data:component-update-manifest stores an opaque boundary continuation and revision
nested_boundaries:
  - track parent identity and stable child ordering
  - unchanged nested boundaries may be omitted independently when the protocol can preserve their DOM
  - otherwise replace the nearest safe changed ancestor
constraints:
  - repeated component calls require stable explicit keys through rule:component-instance-identity
  - browser runtime cannot instantiate undeclared components or select arbitrary server arguments
  - boundary rerendering is side-effect-free and safe to repeat
acceptance:
  - update flag is opt-in and leaves ordinary component output compatible
  - generated route layouts participate automatically without making their parameters client-mutable
  - browser can locate every returned operation target without inspecting component arguments
open_questions:
  - update flag syntax and default boundary element or marker form
  - whether cache-enabled components become update-enabled by shorthand
  - nested boundary preservation versus ancestor replacement in the first milestone
```
