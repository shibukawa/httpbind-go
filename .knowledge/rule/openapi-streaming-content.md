---
id: rule:openapi-streaming-content
type: rule
title: OpenAPI Streaming Response Content Types
---
Stream[T] handlers may document multiple response content types; runtime still negotiates the transport.

```yaml
return_type: "httpbind.Stream[Event]"
possible_content_types:
  - text/event-stream
  - application/x-ndjson
  - application/json
meanings:
  text/event-stream: SSE frames
  application/x-ndjson: NDJSON / JSONL line stream
  application/json: JSON array document of event objects
negotiation: rule:stream-content-negotiation
runtime: content negotiation selects actual format
related:
  - concept:streaming
  - concept:openapi-generation
  - api:new-stream
  - api:stream-write
```
