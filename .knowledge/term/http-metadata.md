---
id: term:http-metadata
type: term
title: HTTP Metadata Tags
---
HTTP-specific request values that must be declared explicitly with path, header, cookie, or method tags.

```yaml
tags:
  path: path segment parameters
  header: HTTP headers
  cookie: cookies
  method: expected HTTP method
must_be_explicit: true
example: |
  type Request struct {
      ID      string `path:"id"`
      Token   string `header:"Authorization"`
      Session string `cookie:"session"`
      Method  string `method:"POST"`
  }
related:
  - concept:request-binding
  - rule:default-input-tag
```
