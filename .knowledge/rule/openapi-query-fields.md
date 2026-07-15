---
id: rule:openapi-query-fields
type: rule
title: OpenAPI Mapping for query Fields
---
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
