---
id: flow:boundary-parameter-update
type: flow
title: Boundary Parameter Update Flow
---
Direct component-subtree update initiated by a declared browser parameter binding.

```yaml
flow:
  trigger: application calls api:client-component-update with boundary handle and parameter patch
  steps:
    - id: collect
      action: client runtime resolves the handle and batches allowed parameter changes
    - id: request
      output: data:boundary-parameter-update-request with instance, revision, continuation, and subtree validators
    - id: reconstruct
      action: server authenticates and reconstructs immutable target inputs from boundary continuation
    - id: validate
      action: decode changes through generated typed parameter validators
    - id: execute
      action: render target requirement:partial-update-boundaries subtree only
    - id: compare
      action: reuse requirement:component-output-cache and omit unchanged descendant validators
    - id: respond
      output: data:component-delta-response scoped to the target subtree
    - id: apply
      action: client accepts latest matching revision, applies operations, and stores next subtree manifest
  failure:
    stale_or_invalid: discard, refresh enclosing state, or perform complete navigation by policy
```
