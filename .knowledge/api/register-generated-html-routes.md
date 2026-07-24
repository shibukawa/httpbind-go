---
id: api:register-generated-html-routes
type: api
title: Register Generated HTML Routes
---
Install all generated filesystem page handlers into an application-owned ServeMux with typed dependencies and policies.

```yaml
source: requirement:generated-html-route-handlers
conceptual_signature: RegisterRoutes(mux *http.ServeMux, deps data:html-route-dependencies, options...) error
behavior:
  - validate mux, dependency groups, generated route conflicts, and required security providers
  - register one generated handler for each data:html-render-route-plan
  - register generated component update and partial navigation endpoints when enabled
  - return startup error before serving when configuration is incomplete
options:
  - handler middleware or wrapper hook
  - authentication and authorization context hook
  - error and invalid-parameter mapping
  - policy:html-update-csrf-protection provider and origin configuration
  - runtime asset and document bootstrap configuration
constraints:
  - registration is an explicit startup call, not package init side effect
  - application owns http.Server, ServeMux, middleware order, lifecycle, and observability
  - manual handlers may coexist when patterns do not conflict
```
