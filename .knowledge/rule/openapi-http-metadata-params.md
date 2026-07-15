---
id: rule:openapi-http-metadata-params
type: rule
title: OpenAPI Parameters for HTTP Metadata
---
path, header, and cookie fields become the corresponding OpenAPI Parameter objects.

```yaml
supported_tags:
  path: path parameter
  header: header parameter
  cookie: cookie parameter
example_request: |
  type CreateUserRequest struct {
      Name  string
      Email string
      OrgID string `path:"org_id"`
  }
generated_request_surface:
  - path parameters
  - query parameters
  - request body
  - validation metadata
related:
  - term:http-metadata
  - concept:openapi-generation
  - concept:request-binding
```
