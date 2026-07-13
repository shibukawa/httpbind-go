# Handler Discovery

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `concept:handler-discovery` | `concept` | Handler Discovery |
| `decision:stdlib-servemux` | `decision` | Stdlib ServeMux Routing |
| `rule:custom-middleware-unwrap` | `rule` | Custom Middleware Unwrap |
| `rule:error-response-discovery` | `rule` | Error Response Discovery |
| `rule:request-model-discovery` | `rule` | Request Model Discovery |
| `rule:response-model-discovery` | `rule` | Response Model Discovery |
| `rule:unsupported-route-patterns` | `rule` | Unsupported Route Patterns |
| `concept:code-generation` | `concept` | Generated Runtime Code |
| `concept:handler-forms` | `concept` | Handler Forms |
| `concept:net-http-handler` | `concept` | Standard net/http Handler |
| `concept:openapi-generation` | `concept` | OpenAPI Generation |
| `concept:route-discovery` | `concept` | Route Discovery |

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

## concept:code-generation

Generator emits bind and write functions, validation, OpenAPI schemas, and streaming metadata without runtime reflection.

```yaml
artifacts:
  - request binding functions
  - response write functions
  - stream write functions
  - validation
  - OpenAPI schemas
  - streaming metadata
function_examples:
  - "func bindCreateUserRequest(r *http.Request) (CreateUserRequest, error)"
  - "func writeCreateUserResponse(w http.ResponseWriter, r *http.Request, resp CreateUserResponse) error"
  - "func writeChatEventStream(w http.ResponseWriter, r *http.Request, stream httpbinder.Stream[ChatEvent]) error"
public_wrappers:
  - api:bind
  - api:write
  - api:write-error
discovery:
  - concept:handler-discovery
  - flow:handler-parse
  - rule:request-model-discovery
  - rule:response-model-discovery
  - rule:error-response-discovery
runtime: no reflection
related:
  - flow:code-generation
  - decision:reflection-free
  - concept:request-binding
  - concept:response-binding
  - concept:openapi-generation
  - concept:streaming
  - concept:stdlib-wrapper-unwrap
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

## concept:route-discovery

Generator discovers static net/http route registrations with literal method and path patterns.

```yaml
supported_basic_patterns:
  - 'mux.HandleFunc("POST /users/{id}", createUserHandler)'
  - 'mux.Handle("POST /users/{id}", http.HandlerFunc(createUserHandler))'
  - 'http.HandleFunc("GET /health", healthHandler)'
  - 'http.Handle("GET /users/{id}", http.HandlerFunc(getUserHandler))'
pattern_requirements:
  - static string route pattern preferred
  - Go 1.22+ method-path form supported via decision:stdlib-servemux
next:
  - resolve registered handler via concept:handler-forms
  - unwrap wrappers via concept:stdlib-wrapper-unwrap
related:
  - concept:handler-discovery
  - rule:unsupported-route-patterns
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
