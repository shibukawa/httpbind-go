# RFC 9457 Problem Details Errors

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `policy:problem-details` | `policy` | RFC 9457 Problem Details Errors |
| `data:problem` | `data` | Problem |
| `rule:error-cause-preservation` | `rule` | Error Cause Preservation |
| `rule:openapi-error-statuses` | `rule` | OpenAPI Statuses from Error Helpers |
| `rule:standard-error-mapping` | `rule` | Standard Error Mapping |
| `system:httpbinder` | `system` | httpbinder Library |
| `api:write-error` | `api` | httpbinder.WriteError |
| `concept:error-helpers` | `concept` | Error Helpers |
| `concept:openapi-generation` | `concept` | OpenAPI Generation |
| `concept:response-binding` | `concept` | Response Binding |

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

## rule:openapi-error-statuses

Recognized httpbinder error constructors add corresponding OpenAPI error responses automatically.

```yaml
constructors_to_status:
  httpbinder.BadRequest: 400
  httpbinder.Unauthorized: 401
  httpbinder.Forbidden: 403
  httpbinder.NotFound: 404
  httpbinder.Conflict: 409
  httpbinder.Internal: 500
  httpbinder.Validation: 400
media_type: application/problem+json
schema: policy:problem-details
discovery: rule:error-response-discovery
related:
  - concept:error-helpers
  - concept:openapi-generation
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

## system:httpbinder

Go library and code generator for typed HTTP request binding, response writing, validation, streaming, and OpenAPI.

```yaml
role: runtime plus ahead-of-time generator
runtime_style: generated code only; no reflection
public_api:
  - api:bind
  - api:write
  - api:write-error
  - concept:error-helpers
primary_inputs:
  - developer-defined Go types
  - struct field tags for source restriction
  - same-package net/http handlers
outputs:
  - generated bind and write functions
  - validation code
  - OpenAPI schemas
  - streaming metadata
related:
  - vision:httpbinder
  - flow:code-generation
  - flow:handler-request
  - concept:request-binding
  - concept:response-binding
  - concept:streaming
  - concept:net-http-handler
  - concept:handler-discovery
  - flow:handler-parse
  - concept:stdlib-wrapper-unwrap
  - policy:problem-details
  - decision:stdlib-servemux
  - concept:openapi-generation
  - concept:openapi-embed
  - api:openapi-json
  - api:openapi-yaml
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

## concept:openapi-generation

OpenAPI is generated from Go source via the shared IR; the document is never hand-edited and stays synchronized with binders, writers, and validation.

```yaml
primary_source: Go source code
openapi_role: derived artifact only
never: manual OpenAPI edits as source of truth
version: decision:openapi-31
schemas_for:
  - request models
  - response models
  - validation constraints
  - error models
  - streaming event models
field_rules:
  - rule:openapi-input-fields
  - rule:openapi-payload-fields
  - rule:openapi-query-fields
  - rule:openapi-http-metadata-params
responses:
  - rule:openapi-success-response
  - rule:openapi-streaming-content
  - rule:openapi-error-statuses
validation_tags: rule:openapi-validation-metadata
errors: policy:problem-details
route_analysis:
  - concept:route-discovery
  - concept:stdlib-wrapper-unwrap
  - flow:handler-parse
artifacts:
  - concept:openapi-embed
  - api:openapi-json
  - api:openapi-yaml
  - concept:openapi-ui
goals: requirement:openapi-goals
pipeline:
  - Go source
  - intermediate representation
  - request binder
  - response writer
  - validation
  - error mapping
  - OpenAPI
related:
  - decision:single-source-of-truth
  - flow:code-generation
  - concept:code-generation
  - concept:handler-discovery
```

## concept:response-binding

Service functions return ordinary Go values; handlers serialize them with api:write.

```yaml
public_api: api:write
service_shape: concept:service-layer
handler_shape: concept:net-http-handler
generated_examples:
  - "func writeCreateUserResponse(w http.ResponseWriter, r *http.Request, resp CreateUserResponse) error"
  - "func writeChatEventStream(w http.ResponseWriter, r *http.Request, stream httpbinder.Stream[ChatEvent]) error"
behavior:
  - serialize success value via api:write
  - map errors via api:write-error and policy:problem-details
related:
  - concept:code-generation
  - concept:streaming
  - system:httpbinder
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
