---
id: rule:response-model-discovery
type: rule
title: Response Model Discovery
---
Response models are discovered from the generic type argument of httpbinder.Write[T](w, r, value), including Stream[T].

```yaml
detection_calls:
  - "httpbinder.Write[T](w, r, value)"
  - "httpbinder.WriteStatus[T](w, r, status, value)"
  - "httpbinder.NewStream[T](w, r)"
ordinary_example: "httpbinder.Write[CreateUserResponse](w, r, output)"
status_example: "httpbinder.WriteStatus[CreateUserResponse](w, r, http.StatusCreated, output)"
streaming_example: |
  stream, err := httpbinder.NewStream[ChatEvent](w, r)
  _ = stream.Write(ChatEvent{...})
model_source: generic type argument T
streaming_type: "httpbinder.Stream[EventType] via NewStream[EventType]"
symbol_identity: rule:go-types-symbol-identity
must_be:
  - github.com/shibukawa/httpbind-go.Write
  - github.com/shibukawa/httpbind-go.WriteStatus
  - github.com/shibukawa/httpbind-go.NewStream
reject: same-named Write/NewStream from other packages
alias_ok: true
openapi_status: rule:openapi-success-status
related:
  - api:write
  - api:write-status
  - api:new-stream
  - concept:response-binding
  - concept:streaming
  - concept:handler-discovery
  - concept:openapi-generation
  - rule:go-types-symbol-identity
  - rule:openapi-success-status
```


