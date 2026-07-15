---
id: concept:streaming
type: concept
title: Typed Streaming
---
Streaming uses a writable Stream[T] obtained from NewStream; handlers call Write repeatedly for each event.

```yaml
api: api:new-stream
type: "httpbinder.Stream[T]"
not:
  - WriteNDJSON batch helper
  - WriteSSE batch helper
handler_shape: |
  stream, err := httpbinder.NewStream[ChatEvent](w, r)
  if err != nil { ... }
  defer stream.Close()
  _ = stream.Write(ChatEvent{Type: "delta", Delta: "hi"})
  _ = stream.Write(ChatEvent{Type: "done"})
service_note: |
  Preferred handler-side API is NewStream + Write.
  Returning Stream[T] from a pure service function remains a future convenience
  if generation wires it to the same runtime writer.
formats:
  - name: sse
    media_type: text/event-stream
    note: Server-Sent Events; data: <json> frames
  - name: ndjson
    media_type: application/x-ndjson
    aliases: [JSONL, application/jsonl, application/ndjson]
    note: one JSON object per line; NOT a single JSON array document
  - name: json-array
    media_type: application/json
    framing: "[obj1,obj2,...]"
    note: single JSON array document; Close writes trailing bracket
    not: JSONL
selection: rule:stream-content-negotiation
openapi: rule:openapi-streaming-content
related:
  - api:new-stream
  - api:stream-write
  - rule:stream-content-negotiation
  - concept:net-http-handler
  - system:httpbinder
```
