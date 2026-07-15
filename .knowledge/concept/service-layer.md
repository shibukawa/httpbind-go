---
id: concept:service-layer
type: concept
title: Service Layer
---
Business logic is ordinary Go functions taking context and typed request, returning typed response and error.

```yaml
shape: |
  func createUser(
      ctx context.Context,
      req CreateUserRequest,
  ) (CreateUserResponse, error)
example_return: |
  return CreateUserResponse{
      ID:    "user_123",
      Name:  req.Name,
      Email: req.Email,
  }, nil
errors:
  - return concept:error-helpers constructors
  - preserve causes with rule:error-cause-preservation
related:
  - concept:net-http-handler
  - concept:response-binding
```
