---
id: rule:openapi-nested-schemas
type: rule
title: OpenAPI Mapping for Nested Request Types
---
Nested structs become component schemas or inline objects; slices become arrays; maps become additionalProperties objects.

```yaml
status: planned
from: rule:nested-request-binding
mapping:
  nested_struct:
    named: components.schemas ref when reused
    anonymous: inline object schema
  slice:
    type: array
    items: element schema
  map_string_T:
    type: object
    additionalProperties: value schema for T
  map_string_any:
    type: object
    additionalProperties: true
  data:file:
    type: string
    format: binary
depth: recursive until scalars or file
related:
  - rule:nested-request-binding
  - concept:openapi-generation
  - rule:openapi-payload-fields
  - data:nested-order-request
```
