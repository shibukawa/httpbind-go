---
id: rule:error-cause-preservation
type: rule
title: Error Cause Preservation
---
HTTP error helpers may wrap an original cause; callers can still use errors.Is, errors.As, and errors.Unwrap.

```yaml
rule: preserve original cause under HTTP error wrapper
example: |
  user, err := repository.Find(ctx, id)
  if err != nil {
      if errors.Is(err, sql.ErrNoRows) {
          return UserResponse{}, httpbind.NotFound(
              Problem{
                  Code:    "user_not_found",
                  Message: "user not found",
              },
              err,
          )
      }
      return UserResponse{}, httpbind.Internal(err)
  }
compatible_apis:
  - errors.Is
  - errors.As
  - errors.Unwrap
client_visibility:
  - public body uses policy:problem-details
  - internal cause logged by api:write-error; not exposed to clients
related:
  - concept:error-helpers
  - data:problem
  - api:write-error
```
