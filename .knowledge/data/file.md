---
id: data:file
type: data
title: httpbinder.File
---
File upload type that is payload-only and binds automatically from multipart/form-data.

```yaml
type: httpbinder.File
source: payload only
media_type: multipart/form-data
example_model: data:upload-avatar-request
example: |
  type UploadAvatarRequest struct {
      UserID string          `path:"user_id"`
      Image  httpbinder.File `payload:"image"`
  }
related:
  - term:payload
  - concept:request-binding
  - data:upload-avatar-request
```
