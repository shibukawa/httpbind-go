# OpenAPI Generation

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `concept:openapi-generation` | `concept` | OpenAPI Generation |
| `decision:openapi-31` | `decision` | OpenAPI 3.1 Target |
| `decision:single-source-of-truth` | `decision` | Single Source of Truth |
| `policy:problem-details` | `policy` | RFC 9457 Problem Details Errors |
| `rule:openapi-error-statuses` | `rule` | OpenAPI Statuses from Error Helpers |
| `rule:openapi-http-metadata-params` | `rule` | OpenAPI Parameters for HTTP Metadata |
| `rule:openapi-input-fields` | `rule` | OpenAPI Mapping for input Fields |
| `rule:openapi-payload-fields` | `rule` | OpenAPI Mapping for payload Fields |
| `rule:openapi-query-fields` | `rule` | OpenAPI Mapping for query Fields |
| `rule:openapi-streaming-content` | `rule` | OpenAPI Streaming Response Content Types |
| `rule:openapi-success-response` | `rule` | OpenAPI Success Response from Write |
| `rule:openapi-validation-metadata` | `rule` | OpenAPI Validation Metadata from Struct Tags |

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

## decision:openapi-31

Generated OpenAPI documents target OpenAPI 3.1; JSON Schema follows the OpenAPI 3.1 rules.

```yaml
status: accepted
version: "3.1"
json_schema: OpenAPI 3.1 dialect
related:
  - concept:openapi-generation
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

## rule:openapi-http-metadata-params

path, header, and cookie fields become the corresponding OpenAPI Parameter objects.

```yaml
supported_tags:
  path: path parameter
  header: header parameter
  cookie: cookie parameter
example_request: |
  type CreateUserRequest struct {
      Name  string
      Email string
      OrgID string `path:"org_id"`
  }
generated_request_surface:
  - path parameters
  - query parameters
  - request body
  - validation metadata
related:
  - term:http-metadata
  - concept:openapi-generation
  - concept:request-binding
```

## rule:openapi-input-fields

Fields with input (or default untagged input) appear as query parameters and in all supported request body media types.

```yaml
tag: input
also: rule:default-input-tag
openapi:
  - query parameter
  - request body application/json
  - request body application/x-www-form-urlencoded
  - request body multipart/form-data
example:
  field: 'Name string `input:"name"`'
  produces:
    - query name
    - json body property name
    - form body property name
    - multipart field name
related:
  - term:input
  - concept:openapi-generation
  - concept:request-binding
```

## rule:openapi-payload-fields

Fields tagged payload appear only in request body media types, not as query parameters.

```yaml
tag: payload
example: 'Name string `payload:"name"`'
openapi_media_types:
  - application/json
  - application/x-www-form-urlencoded
  - multipart/form-data
not_in:
  - query parameters
related:
  - term:payload
  - concept:openapi-generation
```

## rule:openapi-query-fields

Fields tagged query generate only OpenAPI query parameters.

```yaml
tag: query
example: 'Page int `query:"page"`'
openapi:
  - query parameter only
not_in:
  - request body
related:
  - term:query
  - concept:openapi-generation
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

## rule:openapi-success-response

httpbinder.Write[T] discovery generates a 200 OK response with the schema for T.

```yaml
detection: rule:response-model-discovery
call: "httpbinder.Write[UserResponse](...)"
openapi:
  status: 200
  description: OK
  schema: UserResponse
related:
  - api:write
  - concept:response-binding
  - concept:openapi-generation
```

## rule:openapi-validation-metadata

Validation and documentation metadata for OpenAPI schemas is generated from struct tags.

```yaml
supported_metadata:
  - required
  - default
  - enum
  - minimum
  - maximum
  - pattern
  - format
  - deprecated
  - example
  - description
source: Go struct tags
related:
  - concept:openapi-generation
  - concept:request-binding
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
