# Handler Parse Flow

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `flow:handler-parse` | `flow` | Handler Parse Flow |
| `rule:custom-middleware-unwrap` | `rule` | Custom Middleware Unwrap |
| `rule:error-response-discovery` | `rule` | Error Response Discovery |
| `rule:nested-wrapper-unwrap` | `rule` | Nested Wrapper Unwrap |
| `rule:request-model-discovery` | `rule` | Request Model Discovery |
| `rule:response-model-discovery` | `rule` | Response Model Discovery |
| `rule:unsupported-route-patterns` | `rule` | Unsupported Route Patterns |
| `api:bind` | `api` | httpbinder.Bind |
| `api:write` | `api` | httpbinder.Write |
| `concept:error-helpers` | `concept` | Error Helpers |
| `concept:handler-discovery` | `concept` | Handler Discovery |
| `concept:handler-forms` | `concept` | Handler Forms |

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

## rule:custom-middleware-unwrap

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

## rule:nested-wrapper-unwrap

Nested stdlib wrappers are unwrapped outer-to-inner; metadata from each layer is collected.

```yaml
example: |
  mux.Handle(
      "POST /upload",
      http.TimeoutHandler(
          http.MaxBytesHandler(
              http.HandlerFunc(uploadHandler),
              10<<20,
          ),
          30*time.Second,
          "timeout",
      ),
  )
parse_order:
  - unwrap TimeoutHandler
  - unwrap MaxBytesHandler
  - unwrap HandlerFunc
  - analyze uploadHandler
collected_metadata_example:
  timeout: 30s
  max_request_body_bytes: 10485760
related:
  - concept:stdlib-wrapper-unwrap
  - rule:wrapper-timeout-handler
  - rule:wrapper-max-bytes-handler
```

## rule:request-model-discovery

Request models are discovered from the generic type argument of httpbinder.Bind[T](r).

```yaml
detection_call: "httpbinder.Bind[T](r)"
example: "input, err := httpbinder.Bind[CreateUserRequest](r)"
model_source: generic type argument T
related:
  - api:bind
  - concept:request-binding
  - concept:handler-discovery
  - concept:openapi-generation
```

## rule:response-model-discovery

Response models are discovered from the generic type argument of httpbinder.Write[T](w, r, value), including Stream[T].

```yaml
detection_calls:
  - "httpbinder.Write[T](w, r, value)"
  - "httpbinder.NewStream[T](w, r)"
ordinary_example: "httpbinder.Write[CreateUserResponse](w, r, output)"
streaming_example: |
  stream, err := httpbinder.NewStream[ChatEvent](w, r)
  _ = stream.Write(ChatEvent{...})
model_source: generic type argument T
streaming_type: "httpbinder.Stream[EventType] via NewStream[EventType]"
related:
  - api:write
  - api:new-stream
  - concept:response-binding
  - concept:streaming
  - concept:handler-discovery
  - concept:openapi-generation
```

## rule:unsupported-route-patterns

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
escape_hatches:
  - explicit route annotation
  - httpbinder route helper
related:
  - concept:route-discovery
  - concept:handler-discovery
  - decision:stdlib-servemux
```

## api:bind

Generic request binder that maps *http.Request into a typed request struct using generated code.

```yaml
signature: "func Bind[T any](r *http.Request) (T, error)"
example: "input, err := httpbinder.Bind[CreateUserRequest](r)"
behavior:
  - bind query, payload, path, header, cookie, method per field tags
  - return typed value or error
  - no runtime reflection
uses:
  - concept:request-binding
  - concept:code-generation
  - rule:default-input-tag
discovery: rule:request-model-discovery
error_path: api:write-error
related:
  - system:httpbinder
  - concept:net-http-handler
  - concept:handler-discovery
```

## api:write

Generic response writer that serializes a typed value or stream to the HTTP response.

```yaml
signature: "func Write[T any](w http.ResponseWriter, r *http.Request, value T) error"
examples:
  - "httpbinder.Write[CreateUserResponse](w, r, output)"
behavior:
  - serialize ordinary response values
  - no runtime reflection
  - streaming uses api:new-stream (not Write[Stream[T]] for incremental handlers)
uses:
  - concept:response-binding
  - concept:code-generation
discovery: rule:response-model-discovery
related:
  - system:httpbinder
  - concept:net-http-handler
  - concept:handler-discovery
  - api:write-error
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

## concept:handler-discovery

Generator analyzes same-package net/http route registration, handlers, selected stdlib wrappers, and Bind/Write/error call sites.

```yaml
scope: same Go package only
supports:
  - standard net/http handlers
  - http.HandlerFunc
  - handler structs with ServeHTTP
  - selected built-in net/http wrappers
out_of_scope:
  - cross-package handler implementation analysis
pipeline:
  - concept:route-discovery
  - concept:handler-forms
  - concept:stdlib-wrapper-unwrap
  - rule:custom-middleware-unwrap
  - rule:request-model-discovery
  - rule:response-model-discovery
  - rule:error-response-discovery
unsupported: rule:unsupported-route-patterns
related:
  - flow:code-generation
  - concept:code-generation
  - concept:net-http-handler
  - decision:stdlib-servemux
  - concept:openapi-generation
```

## concept:handler-forms

Generator follows same-package handlers as named functions, inline functions, or ServeHTTP structs.

```yaml
forms:
  named_function:
    example: |
      func createUserHandler(w http.ResponseWriter, r *http.Request) {
          input, err := httpbinder.Bind[CreateUserRequest](r)
          httpbinder.Write[CreateUserResponse](w, r, output)
      }
  inline_function:
    example: |
      mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
          input, err := httpbinder.Bind[CreateUserRequest](r)
      })
  handler_struct:
    example: |
      type UserHandler struct{}
      func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
          input, err := httpbinder.Bind[UserRequest](r)
      }
      mux.Handle("POST /users", &UserHandler{})
scope: same package only
related:
  - concept:handler-discovery
  - concept:net-http-handler
  - rule:request-model-discovery
  - rule:response-model-discovery
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
