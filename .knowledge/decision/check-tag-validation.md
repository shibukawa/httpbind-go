---
id: decision:check-tag-validation
type: decision
title: Check Tags as Validation SSOT
---
Adopt a dedicated check struct tag as the single source for runtime validation codegen and OpenAPI constraint metadata.

```yaml
status: accepted
tag_name: check
alternatives_rejected:
  - validate: collides with go-playground/validator mental model
  - binding: collides with gin-style tags
  - rule / constraints: longer, less action-oriented
rationale:
  - short and distinct from popular frameworks
  - aligns with vision:tinybind type-as-SSOT
  - one tag feeds runtime validate and OpenAPI
  - no runtime reflection or runtime tag parsing
out_of_scope_v1:
  - cross-field rules (eqfield)
  - dive into slice element deep validation
  - unique slice items
  - custom i18n messages
  - strict RFC 5322 email
  - uri/url (deferred; ambiguous absolute vs relative)
related:
  - concept:check-validation
  - decision:reflection-free
  - decision:single-source-of-truth
  - concept:openapi-generation
```
