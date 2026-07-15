---
id: rule:unsupported-route-patterns
type: rule
title: Unsupported Route Patterns
---
Dynamic, looped, DI-invoked, and cross-package route registration are intentionally unsupported initially.

```yaml
unsupported_examples:
  - 'mux.HandleFunc("GET " + path, handler)'
  - |
    for _, route := range routes {
        mux.HandleFunc(route.Pattern, route.Handler)
    }
  - fx.Invoke(registerRoutes)
  - someOtherPackage.RegisterRoutes(mux)
escape_hatches:
  - explicit route annotation
  - httpbinder route helper
related:
  - concept:route-discovery
  - concept:handler-discovery
  - decision:stdlib-servemux
```
