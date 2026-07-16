---
id: term:payload
type: term
title: payload Tag
---
Field source restricted to the request body; decoder chosen by Content-Type.

```yaml
tag: 'payload:"name"'
decoder_by_content_type:
  application/json: JSON
  application/x-www-form-urlencoded: form
  multipart/form-data: multipart
json_media_types:
  - application/json
  - text/json
  - '*+json structured syntax suffix (e.g. application/problem+json)'
example: "Name string `payload:\"name\"`"
rest_wire: term:payload-rest
nested: rule:nested-request-binding
openapi: rule:openapi-payload-fields
related:
  - concept:request-binding
  - term:input
  - term:payload-rest
  - data:file
  - rule:openapi-payload-fields
  - rule:payload-rest-map
  - rule:nested-request-binding
```
