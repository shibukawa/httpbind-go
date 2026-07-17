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
    - term:payload-rest
  http_metadata:
    - term:http-metadata
default_field_rule: rule:default-input-tag
validation: concept:check-validation
validation_pipeline: rule:check-codegen-pipeline
files: data:file
nested: rule:nested-request-binding
rest_map: rule:payload-rest-map
examples:
  - data:create-user-request
  - data:search-request
  - data:upload-avatar-request
  - data:patch-with-extras-request
  - data:nested-order-request
payload_media_types:
  - application/json
  - application/x-www-form-urlencoded
  - multipart/form-data
generated_examples:
  - "func bindCreateUserRequest(r *http.Request) (CreateUserRequest, error)"
related:
  - concept:code-generation
  - system:tinybind
  - concept:net-http-handler
  - concept:check-validation
  - term:payload-rest
  - rule:payload-rest-map
  - rule:nested-request-binding
```
