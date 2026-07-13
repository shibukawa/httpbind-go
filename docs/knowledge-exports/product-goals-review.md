# Product Goals

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `requirement:product-goals` | `requirement` | Product Goals |
| `decision:reflection-free` | `decision` | Reflection-Free Runtime |
| `decision:single-source-of-truth` | `decision` | Single Source of Truth |
| `decision:stdlib-servemux` | `decision` | Stdlib ServeMux Routing |
| `policy:problem-details` | `policy` | RFC 9457 Problem Details Errors |
| `api:bind` | `api` | httpbinder.Bind |
| `api:write` | `api` | httpbinder.Write |
| `api:write-error` | `api` | httpbinder.WriteError |
| `concept:net-http-handler` | `concept` | Standard net/http Handler |
| `concept:openapi-generation` | `concept` | OpenAPI Generation |
| `concept:request-binding` | `concept` | Request Binding |
| `concept:streaming` | `concept` | Typed Streaming |

## requirement:product-goals

Core product goals for httpbinder API design and runtime behavior.

```yaml
goals:
  - Go-first API development
  - reflection-free runtime
  - code generation only
  - unified JSON / form / multipart handling
  - browser-friendly APIs
  - curl-friendly APIs
  - RFC 9457 compliant errors
  - automatic OpenAPI generation
  - TinyGo compatible
  - type-safe streaming APIs
maps_to:
  - decision:reflection-free
  - decision:single-source-of-truth
  - decision:stdlib-servemux
  - concept:request-binding
  - concept:streaming
  - concept:net-http-handler
  - policy:problem-details
  - concept:openapi-generation
  - requirement:tinygo-wasm
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

## concept:request-binding

Maps HTTP request values into Go structs via api:bind and generated bind functions.

```yaml
public_api: api:bind
categories:
  user_input:
    - term:input
    - term:query
    - term:payload
  http_metadata:
    - term:http-metadata
default_field_rule: rule:default-input-tag
files: data:file
examples:
  - data:create-user-request
  - data:search-request
  - data:upload-avatar-request
payload_media_types:
  - application/json
  - application/x-www-form-urlencoded
  - multipart/form-data
generated_examples:
  - "func bindCreateUserRequest(r *http.Request) (CreateUserRequest, error)"
related:
  - concept:code-generation
  - system:httpbinder
  - concept:net-http-handler
```

## concept:streaming

Streaming uses a writable Stream[T] obtained from NewStream; handlers call Write repeatedly for each event.

```yaml
api: api:new-stream
type: "httpbinder.Stream[T]"
not:
  - WriteNDJSON batch helper
  - WriteSSE batch helper
handler_shape: |
  stream, err := httpbinder.NewStream[ChatEvent](w, r)
  if err != nil { ... }
  defer stream.Close()
  _ = stream.Write(ChatEvent{Type: "delta", Delta: "hi"})
  _ = stream.Write(ChatEvent{Type: "done"})
service_note: |
  Preferred handler-side API is NewStream + Write.
  Returning Stream[T] from a pure service function remains a future convenience
  if generation wires it to the same runtime writer.
formats:
  - name: sse
    media_type: text/event-stream
    note: Server-Sent Events; data: <json> frames
  - name: ndjson
    media_type: application/x-ndjson
    aliases: [JSONL, application/jsonl, application/ndjson]
    note: one JSON object per line; NOT a single JSON array document
  - name: json-array
    media_type: application/json
    framing: "[obj1,obj2,...]"
    note: single JSON array document; Close writes trailing bracket
    not: JSONL
selection: rule:stream-content-negotiation
openapi: rule:openapi-streaming-content
related:
  - api:new-stream
  - api:stream-write
  - rule:stream-content-negotiation
  - concept:net-http-handler
  - system:httpbinder
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
