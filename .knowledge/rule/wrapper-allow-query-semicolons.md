---
id: rule:wrapper-allow-query-semicolons
type: rule
title: AllowQuerySemicolons Unwrap
---
Unwrap http.AllowQuerySemicolons and analyze the inner handler; optional OpenAPI route metadata only.

```yaml
wrapper: http.AllowQuerySemicolons(h)
example: |
  mux.Handle(
      "GET /search",
      http.AllowQuerySemicolons(
          http.HandlerFunc(searchHandler),
      ),
  )
parse:
  - unwrap h
  - analyze inner handler
openapi:
  schema_change: none by default
  metadata:
    allow_query_semicolons: true
related:
  - concept:stdlib-wrapper-unwrap
  - concept:openapi-generation
```
