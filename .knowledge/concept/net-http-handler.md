---
id: concept:net-http-handler
type: concept
title: Standard net/http Handler
---
Handlers stay standard net/http functions; they bind, call service logic, then Write or WriteError.

```yaml
shape: "func(w http.ResponseWriter, r *http.Request)"
pattern: |
  func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
      input, err := httpbind.Bind[CreateUserRequest](r)
      if err != nil {
          httpbind.WriteError(w, r, err)
          return
      }
      output, err := createUser(r.Context(), input)
      if err != nil {
          httpbind.WriteError(w, r, err)
          return
      }
      httpbind.Write[CreateUserResponse](w, r, output)
  }
layers:
  handler: concept:net-http-handler
  service: concept:service-layer
  bind: api:bind
  write: api:write
  write_error: api:write-error
router: decision:stdlib-servemux
related:
  - system:tinybind
  - concept:handler-discovery
```
