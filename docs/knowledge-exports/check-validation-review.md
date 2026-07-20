# Check Tag Validation

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `concept:check-validation` | `concept` | Check Tag Validation |
| `decision:check-tag-validation` | `decision` | Check Tags as Validation SSOT |
| `decision:reflection-free` | `decision` | Reflection-Free Runtime |
| `decision:single-source-of-truth` | `decision` | Single Source of Truth |
| `rule:check-codegen-pipeline` | `rule` | Check Validation Codegen Pipeline |
| `rule:check-format-validators` | `rule` | Check Format Validators |
| `rule:check-required-semantics` | `rule` | Check Required and Zero-Value Semantics |
| `rule:check-tag-syntax` | `rule` | Check Tag DSL Syntax |
| `rule:check-v1-rule-set` | `rule` | Check Validation v1 Rule Set |
| `rule:openapi-validation-metadata` | `rule` | OpenAPI Validation Metadata from Struct Tags |
| `api:bind` | `api` | httpbinder.Bind |
| `concept:code-generation` | `concept` | Generated Runtime Code |

## concept:check-validation

Struct-tag validation rules on request models generate runtime checks and OpenAPI constraints from one source.

```yaml
status: designed
tag_name: check
intent: replace handwritten Validation/Field checks with generated validate functions
ssot:
  input: Go struct check tags
  outputs:
    - generated validateXxx after bind then defaults
    - OpenAPI required/minimum/maximum/minLength/maxLength/enum/pattern/format/default
pipeline: rule:check-codegen-pipeline
pipeline_order:
  - bind
  - validate
  - apply defaults

syntax: rule:check-tag-syntax
rules: rule:check-v1-rule-set
required_semantics: rule:check-required-semantics
formats: rule:check-format-validators
openapi: rule:openapi-validation-metadata
decision: decision:check-tag-validation
example: |
  type CreateUserRequest struct {
      Name  string `check:"required,minlen=1,maxlen=64"`
      Email string `check:"required,email,maxlen=254"`
      Age   int    `check:"min=0,max=150"`
      ID    string `path:"id" check:"required,uuid"`
  }
related:
  - concept:request-binding
  - concept:code-generation
  - concept:openapi-generation
  - concept:error-helpers
  - api:bind
  - decision:reflection-free
  - decision:single-source-of-truth
  - requirement:tinygo-wasm
  - vision:httpbinder
```

## decision:check-tag-validation

Adopt a dedicated check struct tag as the single source for runtime validation codegen and OpenAPI constraint metadata.

```yaml
status: accepted
tag_name: check
alternatives_rejected:
  - validate: collides with go-playground/validator mental model
  - binding: collides with gin-style tags
  - rule / constraints: longer, less action-oriented
rationale:
  - short and distinct from popular frameworks
  - aligns with vision:httpbinder type-as-SSOT
  - one tag feeds runtime validate and OpenAPI
  - no runtime reflection or runtime tag parsing
out_of_scope_v1:
  - cross-field rules (eqfield)
  - dive into slice element deep validation
  - unique slice items
  - custom i18n messages
  - strict RFC 5322 email
  - uri/url (deferred; ambiguous absolute vs relative)
related:
  - concept:check-validation
  - decision:reflection-free
  - decision:single-source-of-truth
  - concept:openapi-generation
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

## rule:check-codegen-pipeline

Generated bind validates first, then applies defaults; failures become validation errors without runtime tag parsing.

```yaml
order:
  - bind fields from request
  - run validateXxx on bound values
  - apply defaults for still-absent values
  - return value or error
rationale:
  - check before default lets defaults sit outside valid ranges as sentinels
  - example: min=1 with default=-1 distinguishes undefined (becomes -1) from explicit invalid -1 (fails min)
  - default-first would validate the sentinel and reject legitimate absences
validate_presence:
  optional_absent: skip value constraints (min/max/minlen/format/enum/pattern) when field absent
  required_absent: required fails before default
  present_invalid: fail (e.g. explicit -1 with min=1) and never apply default for that field
default_presence:
  query_path_header: apply default only when key was absent and validation passed
  body_json: no presence without pointer; default limited or documented
  note: defaults may be outside check ranges by design (sentinel pattern)
generated_shape: |
  func bindCreateUserRequest(r *http.Request) (CreateUserRequest, error) {
      // bind fields...
      if err := validateCreateUserRequest(&out, presence); err != nil {
          return out, err
      }
      applyDefaultsCreateUserRequest(&out, presence)
      return out, nil
  }
errors:
  style: fixed English templates per rule
  map_to: httpbinder.Validation / Field style problem details
  custom_messages: deferred past v1
  i18n: deferred
messages_examples:
  required: required
  min: "must be >= N"
  minlen: "length must be >= N"
  uuid: "must be a valid uuid"
  enum: "must be one of: a, b"
  date: "must be ISO date"
related:
  - concept:check-validation
  - concept:code-generation
  - api:bind
  - concept:error-helpers
  - rule:standard-error-mapping
  - decision:reflection-free
  - rule:check-required-semantics
  - rule:check-v1-rule-set
```

## rule:check-format-validators

Format shortcuts use pragmatic fixed checks; date/time are ISO-only; email is intentionally non-strict.

```yaml
date:
  accept: "2006-01-02"
  go: time.DateOnly
  openapi_format: date
  reject: non-ISO layouts
time:
  accept: "15:04:05"
  go: time.TimeOnly
  openapi_format: time
datetime:
  tag_name: datetime
  accept: RFC3339
  optional_fallback: RFC3339Nano after RFC3339 fail
  reject: timezone-less local datetimes
  openapi_format: date-time
  naming: tag uses datetime; OpenAPI uses date-time
uuid:
  intent: valid UUID string
  openapi_format: uuid
email:
  strictness: pragmatic not RFC5322
  checks:
    - non-empty when required
    - exactly one '@'
    - non-empty local and domain
    - no whitespace
    - domain contains at least one '.'
  combine_with: maxlen=254 recommended
  openapi_format: email
  escape_hatch: user may add pattern for stricter rules
  avoid: net/mail.ParseAddress as sole check (accepts display-name forms)
pattern:
  engine: Go regexp RE2
  tinygo: regexp supported
  codegen:
    - compile at generation time; invalid pattern is codegen error
    - emit package-level compiled regexp vars
  limits: no backrefs or lookahead; document simple constraints preferred
  syntax: rule:check-tag-syntax
related:
  - concept:check-validation
  - rule:check-v1-rule-set
  - requirement:tinygo-wasm
  - rule:openapi-validation-metadata
```

## rule:check-required-semantics

required must account for Go zero values; path/header absence differs from body/query scalars.

```yaml
problem: Go cannot distinguish omitted vs explicit zero for non-pointer scalars
v1_policy:
  string: required means non-empty
  slice: required means non-empty length
  path_header: missing extraction is required violation
  numeric_bool:
    prefer: pointer types or presence tracking when true required is needed
    v1_safe_default: allow required only on string/slice, or document zero-value pitfalls
  body_json:
    omit vs zero: same without pointer; accept limitation
format_interaction:
  empty_optional_fields: skip email/uuid/date/time/datetime when value empty and not required
  with_required: empty fails required before or instead of format
pipeline_note: rule:check-codegen-pipeline runs validate before default so optional absent skips min/max then may receive out-of-range sentinel defaults
sentinel_example:
  tag: 'check:"min=1,default=-1"'
  absent: after pipeline value is -1 (undefined to app)
  present_minus_one: fails min during validate; default not applied
related:
  - concept:check-validation
  - rule:check-v1-rule-set
  - rule:check-codegen-pipeline
  - concept:request-binding
```

## rule:check-tag-syntax

check tag is a compact CSV-like DSL of rule tokens; enum values use pipe separators; pattern is trailing-only in v1.

```yaml
form: 'check:"rule,rule=value,..."'
token_kinds:
  - bare: required, email, uuid, date, time, datetime
  - key_value: min, max, minlen, maxlen, len, default, enum, pattern
separators:
  rules: ","
  enum_values: "|"
pattern_policy:
  v1: pattern= must be last token in the tag
  reason: commas inside regex break CSV split
  alternatives_deferred:
    - semicolon rule separators
    - quoted pattern values
enum_example: 'check:"required,enum=asc|desc|name"'
pattern_example: 'check:"required,pattern=^[A-Z]{3}$"'
not_compatible_with: go-playground/validator full dialect
parser: codegen only; never interpret tags at runtime
related:
  - concept:check-validation
  - rule:check-v1-rule-set
  - decision:check-tag-validation
  - decision:reflection-free
```

## rule:check-v1-rule-set

v1 check rules cover presence, defaults, numeric bounds, lengths, enums, patterns, and ISO format shortcuts.

```yaml
presence_and_defaults:
  - required
  - default
default_timing: after validate; see rule:check-codegen-pipeline
default_may_be_out_of_range: true
numeric_inclusive:
  - min
  - max
length:
  - minlen
  - maxlen
  - len
set_and_pattern:
  - enum
  - pattern
format_shortcuts:
  - uuid
  - email
  - date
  - time
  - datetime
type_applicability:
  min_max: numeric types only
  minlen_maxlen_len: string and slice
  required: see rule:check-required-semantics
  format_shortcuts: string primarily; skip empty unless required
  enum: comparable scalar or string
deferred_optional:
  - gt
  - gte
  - lt
  - lte
  - uri
  - url
  - format= generic sugar
excluded_v1:
  - eq / ne
  - cross-field
  - dive / element validation beyond outer length
  - unique
  - alpha / alphanum / contains family
  - file size / MIME (separate File rules later)
openapi_map:
  required: required / parameter required
  min: minimum
  max: maximum
  minlen: minLength or minItems
  maxlen: maxLength or maxItems
  len: minLength+maxLength or minItems+maxItems
  enum: enum
  pattern: pattern
  email: format email
  uuid: format uuid
  date: format date
  time: format time
  datetime: format date-time
  default: default
related:
  - concept:check-validation
  - rule:check-tag-syntax
  - rule:check-required-semantics
  - rule:check-format-validators
  - rule:openapi-validation-metadata
```

## rule:openapi-validation-metadata

Validation and documentation metadata for OpenAPI schemas is generated from struct tags, primarily concept:check-validation check tags.

```yaml
primary_source: concept:check-validation
supported_metadata:
  - required
  - default
  - enum
  - minimum
  - maximum
  - minLength
  - maxLength
  - minItems
  - maxItems
  - pattern
  - format
  - deprecated
  - example
  - description
check_mapping: rule:check-v1-rule-set
source: Go struct tags
related:
  - concept:openapi-generation
  - concept:request-binding
  - concept:check-validation
  - rule:check-format-validators
  - decision:check-tag-validation
```

## api:bind

Generic request binder that maps *http.Request into a typed request struct using generated code.

```yaml
signature: "func Bind[T any](r *http.Request) (T, error)"
example: "input, err := httpbinder.Bind[CreateUserRequest](r)"
behavior:
  - bind query, payload, path, header, cookie, method per field tags
  - validate check tags then apply defaults per concept:check-validation
  - return typed value or error
  - no runtime reflection
uses:
  - concept:request-binding
  - concept:code-generation
  - concept:check-validation
  - rule:default-input-tag
  - rule:check-codegen-pipeline
discovery: rule:request-model-discovery
error_path: api:write-error
related:
  - system:httpbinder
  - concept:net-http-handler
  - concept:handler-discovery
  - concept:error-helpers
```

## concept:code-generation

Generator emits bind and write functions, validation, OpenAPI schemas, and streaming metadata without runtime reflection.

```yaml
artifacts:
  - request binding functions
  - response write functions
  - stream write functions
  - validation from concept:check-validation
  - OpenAPI schemas
  - streaming metadata
function_examples:
  - "func bindCreateUserRequest(r *http.Request) (CreateUserRequest, error)"
  - "func validateCreateUserRequest(v *CreateUserRequest) error"
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
  - concept:check-validation
  - rule:check-codegen-pipeline
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
