---
id: policy:json-read-limit
type: policy
title: Bounded JSON Reads
---
JSON document APIs must bound retained input before decoding maps or generated models.

```yaml
status: required
default_bytes: 1048576
coverage:
  - api:decode-json
  - generated HTTP JSON body map reader
controls:
  global: SetMaxJSONBodyBytes
  per_call: DecodeJSONLimit
enforcement:
  - reject known Content-Length above limit before read
  - read at most limit plus one byte when length is unknown
  - never use unbounded io.ReadAll
oversize:
  status: 413
  code: payload_too_large
related:
  - api:decode-json
  - concept:request-binding
  - policy:problem-details
```
