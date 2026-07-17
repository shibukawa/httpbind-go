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
  - cross-package handler leaf
  - Bind/Write type arg to foreign package type (e.g. dto.Request)
  - pointer or complex type arguments unsupported by planner
silent_ignore_today: true
must_become: rule:analysis-diagnostics-check
convention: rule:same-package-convention
escape_hatches:
  - explicit route annotation
  - tinybind route helper
related:
  - concept:route-discovery
  - concept:handler-discovery
  - decision:stdlib-servemux
  - requirement:analysis-diagnostics
  - rule:analysis-diagnostics-check
  - rule:same-package-convention
```

