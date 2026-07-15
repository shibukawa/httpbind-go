---
id: rule:custom-middleware-unwrap
type: rule
title: Custom Middleware Unwrap
---
Custom middleware is best-effort unwrapped when the final handler is statically visible; otherwise require annotation.

```yaml
example: |
  mux.Handle(
      "POST /users",
      Logging(
          Auth(
              http.HandlerFunc(createUserHandler),
          ),
      ),
  )
parse:
  mode: best-effort unwrap
  success_when: inner http.HandlerFunc or handler struct is statically visible
  else: require explicit annotation
related:
  - concept:handler-discovery
  - concept:stdlib-wrapper-unwrap
  - concept:handler-forms
```
