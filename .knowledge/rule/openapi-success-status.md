---
id: rule:openapi-success-status
type: rule
title: OpenAPI Success Status from Write and WriteStatus
---
Success response status codes in OpenAPI come from discovered Write (200) and WriteStatus (literal status when static).

```yaml
status: implemented
api:write:
  openapi_status: 200
  schema: T from Write[T]
api:write-status:
  call: "httpbinder.WriteStatus[T](w, r, status, value)"
  status_arg:
    preferred: integer literal or named const resolvable to int (e.g. http.StatusCreated)
    when_static: emit that status under responses
    when_dynamic: fall back to 200 or document as diagnostic (implementation chooses; prefer 2xx range note)
  schema: T from WriteStatus[T]
  no_content_204:
    may_omit_json_content: true
multiple_writes:
  union_statuses: collect distinct success statuses per operation when multiple Write/WriteStatus sites exist
related:
  - api:write
  - api:write-status
  - rule:openapi-success-response
  - rule:response-model-discovery
  - concept:openapi-generation
```
