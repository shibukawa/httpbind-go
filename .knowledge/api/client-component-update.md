---
id: api:client-component-update
type: api
title: Client Component Update API
---
Request a typed rerender of one explicit component boundary from browser code.

```yaml
source: requirement:boundary-parameter-updates
conceptual_signature: update(boundaryHandle, params, options?) -> Promise<UpdateResult>
preferred_surface: namespaced runtime export rather than a global update symbol
arguments:
  boundaryHandle: stable reference published by data:component-update-manifest
  params: partial object allowed by rule:client-mutable-component-parameters
  options:
    debounce: optional duration or scheduling policy
behavior:
  - resolve current boundary instance, revision, continuation, and subtree validators
  - encode data:boundary-parameter-update-request
  - cancel or supersede older in-flight updates for the same handle
  - apply data:component-delta-response only when revision is current
  - resolve after DOM operations and next manifest state are installed
security:
  - apply policy:html-update-csrf-protection when ambient cookie credentials are used
  - send configured CSRF token in a custom request header, never in target URL
  - target only the generated same-origin update endpoint; attacker-controlled URL input is forbidden
errors:
  unknown_handle: reject without a network request
  invalid_params: reject locally when schema is available and always validate on server
  stale_or_reload: refresh enclosing boundary or navigate by server policy
constraints:
  - API updates boundary state; callers do not rewrite runtime marker attributes directly
  - exact JavaScript namespace and generated typing strategy remain open
```
