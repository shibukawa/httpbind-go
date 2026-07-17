---
id: concept:handler-forms
type: concept
title: Handler Forms
---
Generator follows same-package handlers as named functions, inline functions, or ServeHTTP structs.

```yaml
forms:
  named_function:
    example: |
      func createUserHandler(w http.ResponseWriter, r *http.Request) {
          input, err := httpbind.Bind[CreateUserRequest](r)
          httpbind.Write[CreateUserResponse](w, r, output)
      }
  inline_function:
    example: |
      mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
          input, err := httpbind.Bind[CreateUserRequest](r)
      })
  handler_struct:
    example: |
      type UserHandler struct{}
      func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
          input, err := httpbind.Bind[UserRequest](r)
      }
      mux.Handle("POST /users", &UserHandler{})
scope: same package only
related:
  - concept:handler-discovery
  - concept:net-http-handler
  - rule:request-model-discovery
  - rule:response-model-discovery
```
