---
id: decision:async-boundary-syntax
type: decision
title: Async Fallback Recover Syntax
---
Use one three-state boundary syntax for pending, successful, and failed asynchronous component rendering.

```yaml
source:
  - concept:html-render-runtime-extensions
  - user syntax discussion 2026-07-22
review_gate: proposed syntax requires user approval
preferred_shape: async { primary subtree } fallback { pending subtree } recover(error) { failure subtree }
semantics:
  async: primary subtree may consume requirement:async-external-functions values
  fallback: emitted and flushed while dependencies are pending
  recover: replaces fallback when a dependency returns error, panics, or times out
  error: typed data:async-render-error; raw Go error is unavailable
compiler:
  - async clause introduces a boundary that consumes propagated pending effects
  - fallback and recover clauses must be synchronously renderable
  - recover cannot reference unavailable successful values
  - nested failure is handled by the nearest enclosing matching recover clause
  - expected request cancellation and stale partial-update completion bypass recover
naming:
  benefit: async matches the external modifier; fallback preserves the pending term; recover describes error UI replacement
  caveat: unlike Go recover, this clause handles returned errors and timeouts as well as normalized panics
  alternative: error(error) clause if Go panic association proves misleading
optional_recover:
  proposal: allow omission only when a configured safe default preserves fallback and logs the failure
```
