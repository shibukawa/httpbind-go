# Typed Streaming

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `concept:streaming` | `concept` | Typed Streaming |
| `rule:openapi-streaming-content` | `rule` | OpenAPI Streaming Response Content Types |
| `rule:stream-content-negotiation` | `rule` | Stream Content Negotiation |
| `system:httpbinder` | `system` | httpbinder Library |
| `api:new-stream` | `api` | httpbinder.NewStream |
| `api:stream-write` | `api` | Stream Write and Close |
| `concept:net-http-handler` | `concept` | Standard net/http Handler |
| `concept:code-generation` | `concept` | Generated Runtime Code |
| `concept:response-binding` | `concept` | Response Binding |
| `flow:code-generation` | `flow` | Code Generation Pipeline |

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

## rule:openapi-streaming-content

Stream[T] handlers may document multiple response content types; runtime still negotiates the transport.

```yaml
return_type: "httpbinder.Stream[Event]"
possible_content_types:
  - text/event-stream
  - application/x-ndjson
  - application/json
meanings:
  text/event-stream: SSE frames
  application/x-ndjson: NDJSON / JSONL line stream
  application/json: JSON array document of event objects
negotiation: rule:stream-content-negotiation
runtime: content negotiation selects actual format
related:
  - concept:streaming
  - concept:openapi-generation
  - api:new-stream
  - api:stream-write
```

## rule:stream-content-negotiation

NewStream selects SSE, NDJSON/JSONL, or JSON array from stream query, Accept, User-Agent, then default.

```yaml
formats:
  sse:
    media_type: text/event-stream
    framing: "data: <json>\\n\\n"
  ndjson:
    media_type: application/x-ndjson
    framing: one JSON object per line
    aliases: [JSONL, application/jsonl, application/ndjson]
    not: a JSON array document
  json_array:
    media_type: application/json
    framing: "[obj1,obj2,...]"
    close: writes trailing ] or empty []
    not: JSONL / NDJSON line stream
priority:
  - stream query parameter
  - Accept header
  - User-Agent
  - Default
stream_query:
  values:
    sse: [sse, event-stream, events, eventstream]
    ndjson: [ndjson, jsonl, nd, lines]
    json_array: [json, array, json-array, jsonarray]
accept:
  text/event-stream: sse
  application/x-ndjson: ndjson
  application/ndjson: ndjson
  application/jsonl: ndjson
  application/json: json_array
user_agent_hints:
  browser_like: sse
  curl_like: ndjson
typical:
  browser: Server-Sent Events
  curl: NDJSON (JSONL-style)
  Accept application/json: JSON array document
curl_example:
  command: "curl -N http://localhost/chat"
  note: works without extra Accept; defaults to NDJSON for curl UA when no Accept match
json_array_example:
  command: "curl -N -H 'Accept: application/json' http://localhost/chat"
  note: single JSON array body; distinct from ?stream=jsonl NDJSON lines
default: ndjson
related:
  - concept:streaming
  - api:new-stream
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

## api:new-stream

Creates a typed response stream bound to the ResponseWriter; transport format is chosen from the request.

```yaml
signature: "func NewStream[T any](w http.ResponseWriter, r *http.Request) (*Stream[T], error)"
example: |
  stream, err := httpbinder.NewStream[ChatEvent](w, r)
  if err != nil {
      httpbinder.WriteError(w, r, err)
      return
  }
  defer stream.Close()
behavior:
  - negotiate format via rule:stream-content-negotiation
  - write response headers and status once on first successful open
  - return Stream[T] for repeated api:stream-write
  - flush after each event when ResponseWriter supports Flusher
  - formats: sse | ndjson (JSONL) | json-array
errors:
  - if headers already written
  - if negotiation fails (optional; prefer default)
related:
  - concept:streaming
  - api:stream-write
  - rule:stream-content-negotiation
```

## api:stream-write

Stream[T] supports multiple Write calls for incremental events; Close ends the stream when needed.

```yaml
methods:
  Write:
    signature: "func (s *Stream[T]) Write(v T) error"
    notes:
      - callable many times
      - encodes one event in the negotiated format
      - does not re-send HTTP status/headers
      - for json-array, writes [ on first Write then comma-separated objects
  Close:
    signature: "func (s *Stream[T]) Close() error"
    notes:
      - required for json-array to emit trailing ] (or [] if no Write)
      - optional for NDJSON/SSE line protocols but still recommended via defer
      - idempotent
formats:
  sse: data frames per Write
  ndjson: one JSON object line per Write (JSONL family)
  json-array: single application/json array document across Writes + Close
removed_apis:
  - httpbinder.WriteNDJSON
  - httpbinder.WriteSSE
reason: batch helpers could not be called incrementally without re-writing headers
related:
  - api:new-stream
  - concept:streaming
  - rule:stream-content-negotiation
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

## flow:code-generation

Generator reads same-package handlers and Go types, then emits runtime bind/write functions, validation, OpenAPI, and streaming metadata from one IR.

```yaml
flow:
  trigger: developer defines Go types and net/http handlers
  steps:
    - id: discover-handlers
      action: run flow:handler-parse on same-package registrations
      refs:
        - concept:handler-discovery
        - concept:route-discovery
        - decision:stdlib-servemux
    - id: unwrap-wrappers
      action: unwrap stdlib wrappers and custom middleware when static
      refs:
        - concept:stdlib-wrapper-unwrap
        - rule:nested-wrapper-unwrap
        - rule:custom-middleware-unwrap
    - id: discover-models
      action: detect Bind/Write/error constructors in handler bodies
      refs:
        - rule:request-model-discovery
        - rule:response-model-discovery
        - rule:error-response-discovery
    - id: parse-go-types
      action: analyze discovered struct fields and tags
    - id: build-ir
      action: build shared intermediate representation including route metadata
    - id: emit-binders
      action: generate bind* functions for request types
      refs:
        - concept:request-binding
        - api:bind
    - id: emit-writers
      action: generate write* functions for response and stream types
      refs:
        - concept:response-binding
        - concept:streaming
        - api:write
    - id: emit-validation
      action: generate validation logic
    - id: emit-streaming-metadata
      action: generate streaming transport metadata
      refs:
        - concept:streaming
    - id: emit-openapi
      action: generate OpenAPI 3.1 model, embed, and serve handlers
      refs:
        - concept:openapi-generation
        - concept:openapi-embed
        - api:openapi-json
        - api:openapi-yaml
        - decision:openapi-31
  invariant: all artifacts derive from the same IR
  related:
    - decision:single-source-of-truth
    - system:httpbinder
    - concept:code-generation
    - flow:handler-parse
    - requirement:openapi-goals
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
