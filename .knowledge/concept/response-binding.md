---
id: concept:response-binding
type: concept
title: Response Binding
---
Service functions return ordinary Go values; handlers serialize them with api:write.

```yaml
public_api: api:write
service_shape: concept:service-layer
handler_shape: concept:net-http-handler
generated_examples:
  - "func writeCreateUserResponse(w http.ResponseWriter, r *http.Request, resp CreateUserResponse) error"
  - "func writeChatEventStream(w http.ResponseWriter, r *http.Request, stream httpbinder.Stream[ChatEvent]) error"
behavior:
  - serialize success value via api:write
  - map errors via api:write-error and policy:problem-details
related:
  - concept:code-generation
  - concept:streaming
  - system:httpbinder
```
