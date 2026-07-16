---
id: rule:openapi-payload-rest
type: rule
title: OpenAPI Mapping for payload Rest Fields
---
`payload:"*"` map fields become additionalProperties on the request body object schema.

```yaml
status: implemented
tag: term:payload-rest
openapi:
  in: request body object schema
  mechanism: additionalProperties
  preferred_schema:
    additionalProperties: true
  or_typed:
    additionalProperties:
      # when map[string]json.RawMessage / any
      {}
  not_in:
    - path parameters
    - query parameters
    - header parameters
named_properties: sibling payload/input fields remain normal properties
example: |
  # body with Name + rest
  properties:
    name: { type: string }
  additionalProperties: true
related:
  - term:payload-rest
  - rule:payload-rest-map
  - rule:openapi-payload-fields
  - concept:openapi-generation
```
