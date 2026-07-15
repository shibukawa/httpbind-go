---
id: api:stream-write
type: api
title: Stream Write and Close
---
Stream[T] supports multiple Write calls for incremental events; Close ends the stream when needed.

```yaml
methods:
  Write:
    signature: "func (s *Stream[T]) Write(v T) error"
    notes:
      - callable many times
      - encodes one event in the negotiated format
      - does not re-send HTTP status/headers
      - for json-array, writes [ on first Write then comma-separated objects
  Close:
    signature: "func (s *Stream[T]) Close() error"
    notes:
      - required for json-array to emit trailing ] (or [] if no Write)
      - optional for NDJSON/SSE line protocols but still recommended via defer
      - idempotent
formats:
  sse: data frames per Write
  ndjson: one JSON object line per Write (JSONL family)
  json-array: single application/json array document across Writes + Close
removed_apis:
  - httpbinder.WriteNDJSON
  - httpbinder.WriteSSE
reason: batch helpers could not be called incrementally without re-writing headers
related:
  - api:new-stream
  - concept:streaming
  - rule:stream-content-negotiation
```
