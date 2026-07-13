# Error Helpers

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `concept:error-helpers` | `concept` | Error Helpers |
| `data:problem` | `data` | Problem |
| `policy:problem-details` | `policy` | RFC 9457 Problem Details Errors |
| `rule:error-cause-preservation` | `rule` | Error Cause Preservation |
| `rule:error-response-discovery` | `rule` | Error Response Discovery |
| `rule:standard-error-mapping` | `rule` | Standard Error Mapping |
| `api:write-error` | `api` | httpbinder.WriteError |
| `concept:service-layer` | `concept` | Service Layer |
| `flow:handler-parse` | `flow` | Handler Parse Flow |

## concept:error-helpers

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
```

## data:problem

Application error payload carried by status helpers; includes machine code and human message.

```yaml
fields:
  - name: Code
    type: string
    purpose: machine-readable error code
  - name: Message
    type: string
    purpose: human-readable message
example: |
  Problem{
      Code:    "invalid_email",
      Message: "email is invalid",
  }
used_by:
  - concept:error-helpers
  - policy:problem-details
related:
  - api:write-error
```

## policy:problem-details

Default error format is RFC 9457 Problem Details with field-level locations for validation failures.

```yaml
standard: RFC 9457
default_format: application/problem+json
example:
  type: "..."
  title: Validation failed
  status: 400
  errors:
    - field: email
      location: payload
      message: must be a valid email
error_locations:
  - payload
  - query
  - path
  - header
  - cookie
constructors: concept:error-helpers
problem_type: data:problem
write_path: api:write-error
write_error_behavior:
  - resolve HTTP status
  - convert error to RFC 9457 response
  - log wrapped internal cause
  - hide internal implementation details from clients
mappings: rule:standard-error-mapping
openapi_schema: generated automatically for error responses
openapi_media_type: application/problem+json
related:
  - system:httpbinder
  - concept:response-binding
  - rule:error-cause-preservation
  - concept:openapi-generation
  - rule:openapi-error-statuses
```

## rule:error-cause-preservation

HTTP error helpers may wrap an original cause; callers can still use errors.Is, errors.As, and errors.Unwrap.

```yaml
rule: preserve original cause under HTTP error wrapper
example: |
  user, err := repository.Find(ctx, id)
  if err != nil {
      if errors.Is(err, sql.ErrNoRows) {
          return UserResponse{}, httpbinder.NotFound(
              Problem{
                  Code:    "user_not_found",
                  Message: "user not found",
              },
              err,
          )
      }
      return UserResponse{}, httpbinder.Internal(err)
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

## rule:error-response-discovery

Recognized httpbinder error constructors feed generated OpenAPI error responses.

```yaml
recognized_constructors:
  - httpbinder.BadRequest
  - httpbinder.Unauthorized
  - httpbinder.Forbidden
  - httpbinder.NotFound
  - httpbinder.Conflict
  - httpbinder.Internal
  - httpbinder.Validation
purpose: generate OpenAPI error responses
status_mapping: rule:openapi-error-statuses
media_type: application/problem+json
related:
  - concept:error-helpers
  - policy:problem-details
  - concept:openapi-generation
  - concept:handler-discovery
  - rule:openapi-error-statuses
```

## rule:standard-error-mapping

Built-in mappings convert common Go and parse errors to HTTP status codes; apps may extend or override.

```yaml
mappings:
  - from: JSON parse error
    to: 400 Bad Request
  - from: validation error
    to: 400 Bad Request
  - from: multipart too large
    to: 413 Payload Too Large
  - from: fs.ErrNotExist
    to: 404 Not Found
  - from: context.DeadlineExceeded
    to: 504 Gateway Timeout
  - from: context.Canceled
    to: configurable (e.g. 499)
extensible: true
related:
  - policy:problem-details
  - concept:error-helpers
```

## api:write-error

Writes an error as an RFC 9457 Problem Details response and keeps internal causes out of the client body.

```yaml
signature: "func WriteError(w http.ResponseWriter, r *http.Request, err error)"
example: |
  if err != nil {
      httpbinder.WriteError(w, r, err)
      return
  }
behavior:
  - resolve HTTP status
  - convert error to RFC 9457 response
  - log wrapped internal cause
  - hide internal implementation details from clients
policy: policy:problem-details
helpers: concept:error-helpers
related:
  - system:httpbinder
  - concept:net-http-handler
  - rule:error-cause-preservation
```

## concept:service-layer

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

## flow:handler-parse

Static analysis flow from route registration through wrapper unwrap to Bind, Write, and error discovery.

```yaml
flow:
  trigger: same-package net/http route registration
  steps:
    - id: discover-route
      action: match concept:route-discovery patterns
    - id: unwrap-wrappers
      action: apply concept:stdlib-wrapper-unwrap and rule:nested-wrapper-unwrap
    - id: unwrap-middleware
      action: best-effort rule:custom-middleware-unwrap
    - id: resolve-handler
      action: resolve concept:handler-forms target
    - id: discover-request
      action: rule:request-model-discovery via api:bind
    - id: discover-response
      action: rule:response-model-discovery via api:write
    - id: discover-errors
      action: rule:error-response-discovery via concept:error-helpers
    - id: collect-route-metadata
      action: record wrapper metadata for concept:openapi-generation
  failure:
    unsupported_registration: rule:unsupported-route-patterns
related:
  - concept:handler-discovery
  - flow:code-generation
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
