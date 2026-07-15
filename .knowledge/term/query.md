---
id: term:query
type: term
title: query Tag
---
Field source restricted to URL query parameters only.

```yaml
tag: 'query:"page"'
accepts:
  - URL query parameters
rejects:
  - request body
example: "Page int `query:\"page\"`"
openapi: rule:openapi-query-fields
related:
  - concept:request-binding
  - term:input
  - rule:openapi-query-fields
```
