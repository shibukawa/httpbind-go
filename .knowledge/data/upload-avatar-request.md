---
id: data:upload-avatar-request
type: data
title: UploadAvatarRequest Example
---
Example file upload request with path user id and multipart image payload.

```yaml
type: |
  type UploadAvatarRequest struct {
      UserID string          `path:"user_id"`
      Image  httpbinder.File `payload:"image"`
  }
binding:
  Image: multipart/form-data via data:file
related:
  - data:file
  - term:payload
  - term:http-metadata
  - concept:request-binding
```
