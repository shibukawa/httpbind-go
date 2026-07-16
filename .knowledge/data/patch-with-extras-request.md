---
id: data:patch-with-extras-request
type: data
title: PatchWithExtrasRequest Example
---
Example request with named payload fields plus term:payload-rest for unmapped body keys.

```yaml
type: |
  type PatchWithExtrasRequest struct {
      Name  string         `payload:"name"`
      Email string         `payload:"email"`
      Extra map[string]any `payload:"*"`
  }
binding:
  Name: payload name
  Email: payload email
  Extra: rule:payload-rest-map remaining body keys
json_example: |
  {
    "name": "Ada",
    "email": "ada@example.com",
    "role": "admin",
    "meta": { "source": "import" }
  }
result:
  Name: Ada
  Email: ada@example.com
  Extra:
    role: admin
    meta:
      source: import
note: name and email are not duplicated inside Extra
related:
  - term:payload-rest
  - rule:payload-rest-map
  - term:payload
  - concept:request-binding
```
