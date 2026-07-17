---
id: concept:response-binding
type: concept
title: Response Binding
---
Service functions return ordinary Go values; handlers serialize them with api:write or api:write-status.

```yaml
public_api:
  - api:write
  - api:write-status
service_shape: concept:service-layer
handler_shape: concept:net-http-handler
generated_examples:
  - "func writeCreateUserResponse(w http.ResponseWriter, r *http.Request, resp CreateUserResponse) error"
  - "func writeChatEventStream(w http.ResponseWriter, r *http.Request, stream httpbinder.Stream[ChatEvent]) error"
behavior:
  - serialize success value via api:write (200) or api:write-status (explicit status)
  - map errors via api:write-error and policy:problem-details
openapi: rule:openapi-success-status
related:
  - concept:code-generation
  - concept:streaming
  - system:httpbinder
  - api:write-status
  - rule:openapi-success-status
```

