---
id: rule:default-input-tag
type: rule
title: Default Field Tag is input
---
Untagged struct fields default to input with the field name; tags only restrict value origin.

```yaml
rule:
  when: no struct tag on field
  then: treat as input using field name
equivalence:
  plain: "Name string"
  explicit: "Name string `input:\"name\"`"
intent: most application models need no tags
tags_needed_when: restricting origin to query, payload, path, header, cookie, or method
example_model: data:create-user-request
related:
  - term:input
  - term:http-metadata
  - concept:request-binding
  - data:create-user-request
```
