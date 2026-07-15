---
id: rule:standard-error-mapping
type: rule
title: Standard Error Mapping
---
Built-in mappings convert common Go and parse errors to HTTP status codes; apps may extend or override.

```yaml
mappings:
  - from: JSON parse error
    to: 400 Bad Request
  - from: validation error
    to: 400 Bad Request
  - from: multipart too large
    to: 413 Payload Too Large
  - from: fs.ErrNotExist
    to: 404 Not Found
  - from: context.DeadlineExceeded
    to: 504 Gateway Timeout
  - from: context.Canceled
    to: configurable (e.g. 499)
extensible: true
related:
  - policy:problem-details
  - concept:error-helpers
```
