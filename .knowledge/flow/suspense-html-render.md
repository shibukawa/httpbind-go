---
id: flow:suspense-html-render
type: flow
title: Async Boundary Render Flow
---
Progressive request flow for asynchronous HTML boundaries.

```yaml
flow:
  trigger: exported HTML renderer encounters requirement:suspense-html-streaming
  steps:
    - id: schedule
      action: start required requirement:async-external-functions work under request context
    - id: fallback
      action: allocate boundary ID and render fallback into the normal response stream
    - id: flush
      action: flush initial bytes through the active encoding when supported
    - id: resolve
      action: receive typed completion through the response coordinator
    - id: render-outcome
      action: render primary on success or decision:async-boundary-syntax recover subtree with data:async-render-error on failure
    - id: append-update
      action: append identified template content and fixed replacement instruction
    - id: finish
      action: wait until all boundaries complete or request context cancels
  failure:
    before_commit: existing render error
    async_error: replace fallback with checked recover content
    recover_render_error: preserve fallback and apply outer or server policy
    cancellation: emit no recover update
```
