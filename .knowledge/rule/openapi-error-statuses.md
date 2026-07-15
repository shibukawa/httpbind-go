---
id: rule:openapi-error-statuses
type: rule
title: OpenAPI Statuses from Error Helpers
---
Recognized httpbinder error constructors add corresponding OpenAPI error responses automatically.

```yaml
constructors_to_status:
  httpbinder.BadRequest: 400
  httpbinder.Unauthorized: 401
  httpbinder.Forbidden: 403
  httpbinder.NotFound: 404
  httpbinder.Conflict: 409
  httpbinder.Internal: 500
  httpbinder.Validation: 400
media_type: application/problem+json
schema: policy:problem-details
discovery: rule:error-response-discovery
related:
  - concept:error-helpers
  - concept:openapi-generation
```
