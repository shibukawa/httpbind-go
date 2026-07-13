# Standard net/http Handler

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `concept:net-http-handler` | `concept` | Standard net/http Handler |
| `decision:stdlib-servemux` | `decision` | Stdlib ServeMux Routing |
| `system:httpbinder` | `system` | httpbinder Library |
| `api:bind` | `api` | httpbinder.Bind |
| `api:write` | `api` | httpbinder.Write |
| `api:write-error` | `api` | httpbinder.WriteError |
| `concept:handler-discovery` | `concept` | Handler Discovery |
| `concept:service-layer` | `concept` | Service Layer |

## concept:net-http-handler

Handlers stay standard net/http functions; they bind, call service logic, then Write or WriteError.

```yaml
shape: "func(w http.ResponseWriter, r *http.Request)"
pattern: |
  func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
      input, err := httpbinder.Bind[CreateUserRequest](r)
      if err != nil {
          httpbinder.WriteError(w, r, err)
          return
      }
      output, err := createUser(r.Context(), input)
      if err != nil {
          httpbinder.WriteError(w, r, err)
          return
      }
      httpbinder.Write[CreateUserResponse](w, r, output)
  }
layers:
  handler: concept:net-http-handler
  service: concept:service-layer
  bind: api:bind
  write: api:write
  write_error: api:write-error
router: decision:stdlib-servemux
related:
  - system:httpbinder
  - concept:handler-discovery
```

## decision:stdlib-servemux

Route registration uses Go 1.22+ net/http ServeMux; no framework-specific router is required.

```yaml
status: accepted
router: net/http ServeMux
go_version_min: "1.22"
example: |
  mux := http.NewServeMux()
  mux.HandleFunc(
      "POST /orgs/{org_id}/users",
      CreateUserHandler,
  )
path_params:
  - from: route patterns such as {org_id}
    to: path-tagged request fields
forbidden_requirement:
  - framework-specific router
related:
  - concept:net-http-handler
  - term:http-metadata
  - concept:handler-discovery
  - concept:route-discovery
  - rule:unsupported-route-patterns
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

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
