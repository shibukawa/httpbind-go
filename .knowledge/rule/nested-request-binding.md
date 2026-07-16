---
id: rule:nested-request-binding
type: rule
title: Nested Struct Slice Map Binding
---
Request models may nest structs, slices, and maps; generators emit recursive or per-type bind code without runtime reflection.

```yaml
status: implemented
applies_to: concept:request-binding
supported_field_kinds:
  scalars:
    - string
    - int
    - int64
    - bool
    - float64
  special:
    - data:file
  composite:
    - named struct
    - nested anonymous struct
    - slice of scalar
    - slice of struct
    - map[string]scalar
    - map[string]struct
    - map[string]any when used as term:payload-rest or explicit nested free-form object
media_types:
  application/json:
    structs: nested objects
    slices: JSON arrays
    maps: JSON objects with string keys
  application/x-www-form-urlencoded:
    scalars: flat keys
    nested: deferred or dotted-key optional later; v1 nested full support is JSON-first
  multipart/form-data:
    scalars: form values
    files: data:file at any exported field path allowed by planner
    nested_objects: same JSON-first rule when part is JSON; otherwise flat form fields
codegen:
  style: reflection-free generated code
  approach:
    - emit bind helpers per nested named type
    - or inline nested assignment for small anonymous structs
  registry: top-level api:bind types still registered; nested types may be private helpers
pointer_fields:
  status: optional later
  note: prefer value fields in v1 nested support
check_tags:
  status: concept:check-validation applies to nested leaves when supported
  nested_required: required on nested scalar/struct fields follows rule:check-required-semantics
openapi: rule:openapi-nested-schemas
examples:
  - data:nested-order-request
  - data:patch-with-extras-request
related:
  - concept:request-binding
  - concept:code-generation
  - decision:reflection-free
  - term:payload
  - term:payload-rest
  - rule:payload-rest-map
  - data:file
```
