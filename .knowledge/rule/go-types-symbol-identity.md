---
id: rule:go-types-symbol-identity
type: rule
title: go/types Symbol Identity
---
Resolve call sites with go/types so only the listed stdlib and tinybind symbols participate in route and handler discovery.

```yaml
status: implemented
requirement: requirement:strict-symbol-identity
method: go/types type-checked AST (host generator)
identity: package path + object name (and receiver type for methods)
not: bare Ident/Selector name equality alone
symbol_set:
  default: built-in identities below
  custom: exact identities from requirement:configurable-generator-discovery

route_registration_only:
  - net/http.Handle
  - net/http.HandleFunc
  - "(*net/http.ServeMux).Handle"
  - "(*net/http.ServeMux).HandleFunc"

tinybind_calls_only:
  package: github.com/shibukawa/tinybind-go
  functions:
    - Bind
    - Write
    - WriteStatus
    - WriteError
    - NewStream
    - DecodeJSON
    - EncodeJSON

  error_constructors:
    - BadRequest
    - Unauthorized
    - Forbidden
    - NotFound
    - Conflict
    - PayloadTooLarge
    - Internal
    - Validation
  optional_helpers_if_scanned:
    - Field
    - BindError
    - AsHTTPError

alias_import:
  required: true
  example: |
    import hb "github.com/shibukawa/tinybind-go"
    hb.Bind[T](r)           # recognized via types
    hb.BadRequest(...)      # recognized via types
  forbid: matching only default name "httpbind"

false_positive_reject:
  - otherpkg.HandleFunc(...)
  - otherpkg.Bind[T](...)
  - local type method named Write/Bind
  - mux-like third-party routers named HandleFunc

applies_to:
  - concept:route-discovery
  - rule:request-model-discovery
  - rule:response-model-discovery
  - rule:error-response-discovery
  - DecodeJSON/EncodeJSON type-arg discovery in generator
  - flow:handler-parse
  - requirement:configurable-generator-discovery

implementation_notes:
  - load package with go/packages or equivalent types.Config
  - use types.Info.Uses / Selections to resolve *types.Func
  - compare func.Pkg().Path() and func.Name(); for methods also receiver Named/Pointer
  - keep TinyGo/runtime free of go/types (analysis is host generator only)

related:
  - requirement:strict-symbol-identity
  - concept:handler-discovery
  - concept:route-discovery
  - concept:error-helpers
  - rule:error-response-discovery
  - rule:request-model-discovery
  - rule:response-model-discovery
  - api:bind
  - api:write
  - api:write-error
  - api:new-stream
  - api:decode-json
  - api:encode-json
  - flow:handler-parse
  - concept:code-generation
```
