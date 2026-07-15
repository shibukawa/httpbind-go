---
id: rule:openapi-success-response
type: rule
title: OpenAPI Success Response from Write
---
httpbinder.Write[T] discovery generates a 200 OK response with the schema for T.

```yaml
detection: rule:response-model-discovery
call: "httpbinder.Write[UserResponse](...)"
openapi:
  status: 200
  description: OK
  schema: UserResponse
related:
  - api:write
  - concept:response-binding
  - concept:openapi-generation
```
