---
id: requirement:boundary-parameter-updates
type: requirement
title: Boundary Parameter Updates
---
Rerender one explicit component boundary when a declared browser-controlled parameter changes.

```yaml
source: concept:html-render-runtime-extensions
flow: flow:boundary-parameter-update
boundary: requirement:partial-update-boundaries
trigger:
  primary: api:client-component-update with boundary handle and allowed parameter patch
  declarative_binding: deferred; no generic data-param attribute is reserved
request: data:boundary-parameter-update-request
server_execution:
  - authenticate request and validate the boundary continuation
  - reconstruct immutable component inputs and current mutable state
  - type-check and apply only allowlisted parameter changes
  - execute the generated target component and its descendants without rerunning page, layouts, or ancestors
  - compare descendant validators and return data:component-delta-response
client_execution:
  - batch changes configured for the same event turn
  - cancel superseded requests when possible
  - apply a response only when its base revision matches the latest accepted state
  - store returned boundary revision, continuation, and subtree manifest
fallback:
  stale_or_invalid_capability: refresh page or enclosing navigation boundary
  unsupported_structural_delta: replace target boundary root
  no_javascript: ordinary form submission or complete navigation when provided
compatibility: requirement:html-rendering-compatibility
acceptance:
  - api:client-component-update rerenders only the addressed declared boundary subtree
  - arbitrary DOM edits cannot select a server component or mutate an unlisted parameter
  - rapid out-of-order responses cannot restore an older boundary state
open_questions:
  - generated endpoint routing and HTTP method
  - exact namespaced JavaScript export and boundary-handle syntax
  - stateless signed or encrypted continuation versus server-side instance store
  - history and URL synchronization for local parameters
  - validation error and loading UI integration with requirement:suspense-html-streaming
```
