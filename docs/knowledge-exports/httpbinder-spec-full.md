# httpbinder Vision

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `vision:httpbinder` | `vision` | httpbinder Vision |
| `decision:reflection-free` | `decision` | Reflection-Free Runtime |
| `decision:single-source-of-truth` | `decision` | Single Source of Truth |
| `decision:stdlib-servemux` | `decision` | Stdlib ServeMux Routing |
| `system:httpbinder` | `system` | httpbinder Library |
| `api:bind` | `api` | httpbinder.Bind |
| `api:write` | `api` | httpbinder.Write |
| `api:write-error` | `api` | httpbinder.WriteError |
| `concept:code-generation` | `concept` | Generated Runtime Code |
| `concept:net-http-handler` | `concept` | Standard net/http Handler |
| `concept:openapi-generation` | `concept` | OpenAPI Generation |
| `requirement:tinygo-wasm` | `requirement` | TinyGo and WASM Support |

## vision:httpbinder

httpbinder is a code-generation-first library that bridges Go types and HTTP APIs without runtime reflection.

```yaml
source_of_truth:
  - Go types only
generated_from_types:
  - request binding
  - response serialization
  - streaming responses
  - error handling
  - validation
  - OpenAPI generation
principles:
  - decision:single-source-of-truth
  - decision:reflection-free
targets:
  - system:httpbinder
  - requirement:tinygo-wasm
  - concept:code-generation
  - concept:openapi-generation
  - concept:net-http-handler
  - decision:stdlib-servemux
public_runtime:
  - api:bind
  - api:write
  - api:write-error
```

## decision:reflection-free

Reflection is intentionally unsupported; binding and serialization use generated code only.

```yaml
status: accepted
forbidden:
  - runtime reflection
  - runtime tag parsing
rationale:
  - better performance
  - smaller binaries
  - TinyGo compatibility
  - WASM compatibility
  - compile-time validation
  - OpenAPI consistency with runtime
related:
  - vision:httpbinder
  - requirement:tinygo-wasm
  - concept:code-generation
  - decision:single-source-of-truth
```

## decision:single-source-of-truth

Developers define Go types only; all binders, writers, validation, OpenAPI, and streaming metadata are generated.

```yaml
authority: Go types
not_authority:
  - OpenAPI document as primary input
pipeline:
  - from: Go types
    to: flow:code-generation
  - from: flow:code-generation
    artifacts:
      - request binder
      - response writer
      - validation
      - OpenAPI
      - streaming metadata
runtime: generated code is the implementation
related:
  - vision:httpbinder
  - concept:code-generation
  - concept:openapi-generation
  - concept:openapi-embed
  - requirement:openapi-goals
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

## requirement:tinygo-wasm

TinyGo is a first-class target; reflection-free generated code must work on TinyGo, WebAssembly, and embedded environments.

```yaml
priority: first-class
targets:
  - TinyGo
  - WebAssembly
  - embedded environments
depends_on:
  - decision:reflection-free
related:
  - vision:httpbinder
  - system:httpbinder
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
