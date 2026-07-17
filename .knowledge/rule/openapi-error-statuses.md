---
id: rule:openapi-error-statuses
type: rule
title: OpenAPI Statuses from Error Helpers
---
Recognized tinybind error constructors add corresponding OpenAPI error responses automatically.

```yaml
constructors_to_status:
  httpbind.BadRequest: 400
  httpbind.Unauthorized: 401
  httpbind.Forbidden: 403
  httpbind.NotFound: 404
  httpbind.Conflict: 409
  httpbind.Internal: 500
  httpbind.Validation: 400
media_type: application/problem+json
schema: policy:problem-details
discovery: rule:error-response-discovery
related:
  - concept:error-helpers
  - concept:openapi-generation
```
