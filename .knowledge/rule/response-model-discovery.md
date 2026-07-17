---
id: rule:response-model-discovery
type: rule
title: Response Model Discovery
---
Response models are discovered from the generic type argument of httpbind.Write[T](w, r, value), including Stream[T].

```yaml
detection_calls:
  - "httpbind.Write[T](w, r, value)"
  - "httpbind.WriteStatus[T](w, r, status, value)"
  - "httpbind.NewStream[T](w, r)"
ordinary_example: "httpbind.Write[CreateUserResponse](w, r, output)"
status_example: "httpbind.WriteStatus[CreateUserResponse](w, r, http.StatusCreated, output)"
streaming_example: |
  stream, err := httpbind.NewStream[ChatEvent](w, r)
  _ = stream.Write(ChatEvent{...})
model_source: generic type argument T
streaming_type: "httpbind.Stream[EventType] via NewStream[EventType]"
symbol_identity: rule:go-types-symbol-identity
must_be:
  - github.com/shibukawa/tinybind-go.Write
  - github.com/shibukawa/tinybind-go.WriteStatus
  - github.com/shibukawa/tinybind-go.NewStream
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


