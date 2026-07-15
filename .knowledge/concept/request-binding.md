---
id: concept:request-binding
type: concept
title: Request Binding
---
Maps HTTP request values into Go structs via api:bind and generated bind functions, then applies concept:check-validation.

```yaml
public_api: api:bind
categories:
  user_input:
    - term:input
    - term:query
    - term:payload
  http_metadata:
    - term:http-metadata
default_field_rule: rule:default-input-tag
validation: concept:check-validation
validation_pipeline: rule:check-codegen-pipeline
files: data:file
examples:
  - data:create-user-request
  - data:search-request
  - data:upload-avatar-request
payload_media_types:
  - application/json
  - application/x-www-form-urlencoded
  - multipart/form-data
generated_examples:
  - "func bindCreateUserRequest(r *http.Request) (CreateUserRequest, error)"
related:
  - concept:code-generation
  - system:httpbinder
  - concept:net-http-handler
  - concept:check-validation
```
