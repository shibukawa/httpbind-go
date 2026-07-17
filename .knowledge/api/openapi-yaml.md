---
id: api:openapi-yaml
type: api
title: OpenAPI YAML Handler
---
Generated handler serializes the embedded OpenAPI document as YAML.

```yaml
signature: |
  func OpenAPIYAML(
      w http.ResponseWriter,
      r *http.Request,
  )
recommended_mount: "GET /openapi.yaml"
example: |
  mux.HandleFunc(
      "GET /openapi.yaml",
      httpbind.OpenAPIYAML,
  )
source: concept:openapi-embed
related:
  - api:openapi-json
  - concept:openapi-generation
  - concept:openapi-ui
```
