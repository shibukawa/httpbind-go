---
id: api:openapi-json
type: api
title: OpenAPI JSON Handler
---
Generated handler serializes the embedded OpenAPI document as JSON.

```yaml
signature: |
  func OpenAPIJSON(
      w http.ResponseWriter,
      r *http.Request,
  )
recommended_mount: "GET /openapi.json"
example: |
  mux.HandleFunc(
      "GET /openapi.json",
      httpbinder.OpenAPIJSON,
  )
source: concept:openapi-embed
related:
  - api:openapi-yaml
  - concept:openapi-generation
  - concept:openapi-ui
```
