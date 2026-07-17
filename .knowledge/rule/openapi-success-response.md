---
id: rule:openapi-success-response
type: rule
title: OpenAPI Success Response from Write
---
httpbind.Write[T] discovery generates a 200 OK response with the schema for T; non-200 success uses api:write-status.

```yaml
detection: rule:response-model-discovery
call: "httpbind.Write[UserResponse](...)"
openapi:
  status: 200
  description: OK
  schema: UserResponse
other_success_statuses: rule:openapi-success-status
related:
  - api:write
  - api:write-status
  - rule:openapi-success-status
  - concept:response-binding
  - concept:openapi-generation
```

