---
id: data:file
type: data
title: httpbind.File
---
File upload type that is payload-only and binds automatically from multipart/form-data.

```yaml
type: httpbind.File
source: payload only
media_type: multipart/form-data
example_model: data:upload-avatar-request
example: |
  type UploadAvatarRequest struct {
      UserID string          `path:"user_id"`
      Image  httpbind.File `payload:"image"`
  }
related:
  - term:payload
  - concept:request-binding
  - data:upload-avatar-request
```
