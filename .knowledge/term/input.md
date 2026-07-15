---
id: term:input
type: term
title: input Tag
---
User-provided field source that accepts URL query parameters or request payload.

```yaml
tag: 'input:"name"'
accepts:
  - query string
  - JSON body
  - form body
  - multipart fields
default_untagged: true
example_model: data:create-user-request
accepted_shapes:
  - "POST /users?name=Alice&email=a@example.com"
  - '{"name":"Alice","email":"a@example.com"}'
  - |
    name=Alice
    email=a@example.com
  - multipart form fields with same names
openapi: rule:openapi-input-fields
related:
  - concept:request-binding
  - term:query
  - term:payload
  - rule:default-input-tag
  - rule:openapi-input-fields
  - data:create-user-request
```
