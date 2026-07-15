---
id: data:create-user-request
type: data
title: CreateUserRequest Example
---
Example request model mixing default input fields with explicit path and header metadata.

```yaml
implicit_form: |
  type CreateUserRequest struct {
      Name  string
      Email string
      OrgID string `path:"org_id"`
      Token string `header:"Authorization"`
  }
explicit_equivalent: |
  type CreateUserRequest struct {
      Name  string `input:"name"`
      Email string `input:"email"`
      OrgID string `path:"org_id"`
      Token string `header:"Authorization"`
  }
binding_examples:
  query: "POST /users?name=Alice&email=a@example.com"
  json:
    name: Alice
    email: a@example.com
  form: |
    name=Alice
    email=a@example.com
  multipart: |
    Content-Type: multipart/form-data
    name=Alice
    email=a@example.com
note: all shapes bind to the same Go type via term:input defaults
related:
  - rule:default-input-tag
  - term:input
  - term:http-metadata
  - concept:request-binding
  - api:bind
```
