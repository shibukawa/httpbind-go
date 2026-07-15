---
id: rule:response-model-discovery
type: rule
title: Response Model Discovery
---
Response models are discovered from the generic type argument of httpbinder.Write[T](w, r, value), including Stream[T].

```yaml
detection_calls:
  - "httpbinder.Write[T](w, r, value)"
  - "httpbinder.NewStream[T](w, r)"
ordinary_example: "httpbinder.Write[CreateUserResponse](w, r, output)"
streaming_example: |
  stream, err := httpbinder.NewStream[ChatEvent](w, r)
  _ = stream.Write(ChatEvent{...})
model_source: generic type argument T
streaming_type: "httpbinder.Stream[EventType] via NewStream[EventType]"
related:
  - api:write
  - api:new-stream
  - concept:response-binding
  - concept:streaming
  - concept:handler-discovery
  - concept:openapi-generation
```
