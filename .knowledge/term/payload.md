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
example: "Name string `payload:\"name\"`"
openapi: rule:openapi-payload-fields
related:
  - concept:request-binding
  - term:input
  - data:file
  - rule:openapi-payload-fields
```
