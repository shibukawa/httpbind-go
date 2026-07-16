---
id: term:payload-rest
type: term
title: payload Rest Tag
---
Special payload wire name `*` that captures body object keys not bound to sibling fields.

```yaml
tag: 'payload:"*"'
wire: "*"
field_go_types:
  - map[string]any
  - map[string]json.RawMessage
intent: keep unknown or extra JSON/object keys without listing every field
not_a_key: "*" is not looked up as a body property name
source: term:payload only; not path/header/cookie/query
example: |
  type PatchUserRequest struct {
      Name  string         `payload:"name"`
      Extra map[string]any `payload:"*"`
  }
related:
  - term:payload
  - rule:payload-rest-map
  - data:patch-with-extras-request
  - concept:request-binding
```
