# Request Binding

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `concept:request-binding` | `concept` | Request Binding |
| `data:create-user-request` | `data` | CreateUserRequest Example |
| `data:file` | `data` | httpbinder.File |
| `data:search-request` | `data` | SearchRequest Example |
| `data:upload-avatar-request` | `data` | UploadAvatarRequest Example |
| `rule:default-input-tag` | `rule` | Default Field Tag is input |
| `system:httpbinder` | `system` | httpbinder Library |
| `api:bind` | `api` | httpbinder.Bind |
| `concept:code-generation` | `concept` | Generated Runtime Code |
| `concept:net-http-handler` | `concept` | Standard net/http Handler |
| `term:http-metadata` | `term` | HTTP Metadata Tags |
| `term:input` | `term` | input Tag |

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

## data:create-user-request

Example request model mixing default input fields with explicit path and header metadata.

```yaml
implicit_form: |
  type CreateUserRequest struct {
      Name  string
      Email string
      OrgID string `path:"org_id"`
      Token string `header:"Authorization"`
  }
explicit_equivalent: |
  type CreateUserRequest struct {
      Name  string `input:"name"`
      Email string `input:"email"`
      OrgID string `path:"org_id"`
      Token string `header:"Authorization"`
  }
binding_examples:
  query: "POST /users?name=Alice&email=a@example.com"
  json:
    name: Alice
    email: a@example.com
  form: |
    name=Alice
    email=a@example.com
  multipart: |
    Content-Type: multipart/form-data
    name=Alice
    email=a@example.com
note: all shapes bind to the same Go type via term:input defaults
related:
  - rule:default-input-tag
  - term:input
  - term:http-metadata
  - concept:request-binding
  - api:bind
```

## data:file

File upload type that is payload-only and binds automatically from multipart/form-data.

```yaml
type: httpbinder.File
source: payload only
media_type: multipart/form-data
example_model: data:upload-avatar-request
example: |
  type UploadAvatarRequest struct {
      UserID string          `path:"user_id"`
      Image  httpbinder.File `payload:"image"`
  }
related:
  - term:payload
  - concept:request-binding
  - data:upload-avatar-request
```

## data:search-request

Example request that restricts fields to query-only or payload-only sources.

```yaml
type: |
  type SearchRequest struct {
      Keyword string `query:"keyword"`
      Page    int    `query:"page"`
      Name    string `payload:"name"`
      Email   string `payload:"email"`
  }
payload_formats:
  - application/json
  - application/x-www-form-urlencoded
  - multipart/form-data
related:
  - term:query
  - term:payload
  - concept:request-binding
```

## data:upload-avatar-request

Example file upload request with path user id and multipart image payload.

```yaml
type: |
  type UploadAvatarRequest struct {
      UserID string          `path:"user_id"`
      Image  httpbinder.File `payload:"image"`
  }
binding:
  Image: multipart/form-data via data:file
related:
  - data:file
  - term:payload
  - term:http-metadata
  - concept:request-binding
```

## rule:default-input-tag

Untagged struct fields default to input with the field name; tags only restrict value origin.

```yaml
rule:
  when: no struct tag on field
  then: treat as input using field name
equivalence:
  plain: "Name string"
  explicit: "Name string `input:\"name\"`"
intent: most application models need no tags
tags_needed_when: restricting origin to query, payload, path, header, cookie, or method
example_model: data:create-user-request
related:
  - term:input
  - term:http-metadata
  - concept:request-binding
  - data:create-user-request
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

## term:http-metadata

HTTP-specific request values that must be declared explicitly with path, header, cookie, or method tags.

```yaml
tags:
  path: path segment parameters
  header: HTTP headers
  cookie: cookies
  method: expected HTTP method
must_be_explicit: true
example: |
  type Request struct {
      ID      string `path:"id"`
      Token   string `header:"Authorization"`
      Session string `cookie:"session"`
      Method  string `method:"POST"`
  }
related:
  - concept:request-binding
  - rule:default-input-tag
```

## term:input

User-provided field source that accepts URL query parameters or request payload.

```yaml
tag: 'input:"name"'
accepts:
  - query string
  - JSON body
  - form body
  - multipart fields
default_untagged: true
example_model: data:create-user-request
accepted_shapes:
  - "POST /users?name=Alice&email=a@example.com"
  - '{"name":"Alice","email":"a@example.com"}'
  - |
    name=Alice
    email=a@example.com
  - multipart form fields with same names
openapi: rule:openapi-input-fields
related:
  - concept:request-binding
  - term:query
  - term:payload
  - rule:default-input-tag
  - rule:openapi-input-fields
  - data:create-user-request
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
