---
id: concept:route-discovery
type: concept
title: Route Discovery
---
Generator discovers static net/http route registrations with literal method and path patterns.

```yaml
supported_basic_patterns:
  - 'mux.HandleFunc("POST /users/{id}", createUserHandler)'
  - 'mux.Handle("POST /users/{id}", http.HandlerFunc(createUserHandler))'
  - 'http.HandleFunc("GET /health", healthHandler)'
  - 'http.Handle("GET /users/{id}", http.HandlerFunc(getUserHandler))'
pattern_requirements:
  - static string route pattern preferred
  - Go 1.22+ method-path form supported via decision:stdlib-servemux
next:
  - resolve registered handler via concept:handler-forms
  - unwrap wrappers via concept:stdlib-wrapper-unwrap
related:
  - concept:handler-discovery
  - rule:unsupported-route-patterns
```
