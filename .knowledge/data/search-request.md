---
id: data:search-request
type: data
title: SearchRequest Example
---
Example request that restricts fields to query-only or payload-only sources.

```yaml
type: |
  type SearchRequest struct {
      Keyword string `query:"keyword"`
      Page    int    `query:"page"`
      Name    string `payload:"name"`
      Email   string `payload:"email"`
  }
payload_formats:
  - application/json
  - application/x-www-form-urlencoded
  - multipart/form-data
related:
  - term:query
  - term:payload
  - concept:request-binding
```
