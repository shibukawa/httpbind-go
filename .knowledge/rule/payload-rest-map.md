---
id: rule:payload-rest-map
type: rule
title: Payload Rest Map Binding
---
After named body fields bind, remaining object keys populate the single `payload:"*"` map field.

```yaml
status: implemented
tag: term:payload-rest
algorithm:
  1: decode body object by Content-Type (JSON object required for full semantics)
  2: bind sibling payload/input body fields by wire name as usual
  3: collect keys present in the body object that no sibling field consumed
  4: assign that key/value map to the rest field
exclusion:
  - keys bound by explicit payload/input/json wire names on the same struct
  - path header cookie method fields never consume body keys
cardinality:
  max_rest_fields_per_struct: 1
  zero_ok: true
type_rules:
  required_kind: map with string keys
  preferred:
    - map[string]any
    - map[string]json.RawMessage
  reject: non-map rest field type
value_preservation:
  json: keep decoded JSON values (numbers objects arrays bools strings null)
  form_urlencoded: remaining form keys as string values when present
  multipart: remaining non-file form values as strings; file parts not dumped into rest unless also modeled as data:file
empty:
  no_remaining_keys: empty non-nil map or nil per generator convention; prefer empty non-nil map
errors:
  - multiple payload:"*" fields: generate-time error
  - rest field wrong type: generate-time error
  - body not an object when rest present and JSON: bind error 400
openapi: rule:openapi-payload-rest
example_model: data:patch-with-extras-request
related:
  - term:payload-rest
  - term:payload
  - concept:request-binding
  - concept:code-generation
  - rule:nested-request-binding
```
