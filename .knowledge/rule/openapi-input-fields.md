---
id: rule:openapi-input-fields
type: rule
title: OpenAPI Mapping for input Fields
---
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
