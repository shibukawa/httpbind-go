---
id: rule:openapi-payload-fields
type: rule
title: OpenAPI Mapping for payload Fields
---
Fields tagged payload appear only in request body media types, not as query parameters.

```yaml
tag: payload
example: 'Name string `payload:"name"`'
openapi_media_types:
  - application/json
  - application/x-www-form-urlencoded
  - multipart/form-data
not_in:
  - query parameters
related:
  - term:payload
  - concept:openapi-generation
```
