---
id: concept:openapi-embed
type: concept
title: Embedded OpenAPI Document
---
Generator builds an in-memory OpenAPI model and embeds it in generated Go code for runtime serving.

```yaml
generated_file_example: httpbinder_openapi_gen.go
embedded_symbol_example: "var generatedOpenAPI = ..."
includes:
  - embedded OpenAPI specification
  - JSON serve handler
  - YAML serve handler
  - schema metadata
handlers:
  - api:openapi-json
  - api:openapi-yaml
related:
  - concept:openapi-generation
  - decision:openapi-31
  - requirement:openapi-goals
```
