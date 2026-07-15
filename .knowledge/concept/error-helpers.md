---
id: concept:error-helpers
type: concept
title: Error Helpers
---
Convenience constructors for HTTP errors, validation fields, and optional wrapped causes.

```yaml
status_helpers:
  - httpbinder.BadRequest
  - httpbinder.Unauthorized
  - httpbinder.Forbidden
  - httpbinder.NotFound
  - httpbinder.Conflict
  - httpbinder.Internal
  - httpbinder.Validation
openapi_discovery: rule:error-response-discovery
problem_payload: data:problem
built_in_examples:
  - |
    return UserResponse{}, httpbinder.BadRequest(
        Problem{Code: "invalid_email", Message: "email is invalid"},
    )
  - |
    return UserResponse{}, httpbinder.NotFound(
        Problem{Code: "user_not_found", Message: "user not found"},
    )
  - |
    return UserResponse{}, httpbinder.Conflict(
        Problem{Code: "duplicate_email", Message: "email already exists"},
    )
  - "return UserResponse{}, httpbinder.Internal(err)"
validation_example: |
  return UserResponse{}, httpbinder.Validation(
      httpbinder.Field("email", "payload", "must be a valid email"),
      httpbinder.Field("age", "payload", "must be greater than or equal to 18"),
  )
generated_validation: concept:check-validation
note: handwritten Validation/Field remains for domain rules; field-level input checks move to check tags
field_helper:
  name: httpbinder.Field
  args:
    - field name
    - location (payload|query|path|header|cookie)
    - message
cause_wrapping: rule:error-cause-preservation
response_writer: api:write-error
related:
  - policy:problem-details
  - rule:standard-error-mapping
  - rule:error-response-discovery
  - data:problem
  - concept:check-validation
  - rule:check-codegen-pipeline
```
