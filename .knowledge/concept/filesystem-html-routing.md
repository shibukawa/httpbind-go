---
id: concept:filesystem-html-routing
type: concept
title: Filesystem HTML Routing
---
Opt-in route tree that generates typed page routes, ancestor layout chains, and partial-update reuse boundaries.

```yaml
evidence:
  source: user design discussion
  received: 2026-07-23
review_gate: proposed requirements require user approval
convention: decision:html-route-file-conventions
route_plan: data:html-render-route-plan
requirements:
  - requirement:generated-html-route-handlers
  - requirement:html-runtime-bootstrap
  - requirement:typed-html-route-parameters
  - requirement:layout-chain-discovery
  - requirement:layout-reuse-boundaries
  - requirement:nested-layout-composition
compatibility: requirement:html-rendering-compatibility
go_package: decision:html-route-go-package-model
principles:
  - generate route structure and typed bindings; never parse templates at request time
  - prefer convention defaults with explicit page override
  - scope layout inputs so unaffected ancestor wrappers remain reusable
  - keep document bootstrap ownership separate from the root index page
  - generate handlers and inject typed data dependencies instead of requiring one Go package per route folder
```
