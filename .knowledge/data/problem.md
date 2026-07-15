---
id: data:problem
type: data
title: Problem
---
Application error payload carried by status helpers; includes machine code and human message.

```yaml
fields:
  - name: Code
    type: string
    purpose: machine-readable error code
  - name: Message
    type: string
    purpose: human-readable message
example: |
  Problem{
      Code:    "invalid_email",
      Message: "email is invalid",
  }
used_by:
  - concept:error-helpers
  - policy:problem-details
related:
  - api:write-error
```
