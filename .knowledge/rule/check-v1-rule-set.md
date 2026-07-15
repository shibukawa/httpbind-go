---
id: rule:check-v1-rule-set
type: rule
title: Check Validation v1 Rule Set
---
v1 check rules cover presence, defaults, numeric bounds, lengths, enums, patterns, and ISO format shortcuts.

```yaml
presence_and_defaults:
  - required
  - default
default_timing: after validate; see rule:check-codegen-pipeline
default_may_be_out_of_range: true
numeric_inclusive:
  - min
  - max
length:
  - minlen
  - maxlen
  - len
set_and_pattern:
  - enum
  - pattern
format_shortcuts:
  - uuid
  - email
  - date
  - time
  - datetime
type_applicability:
  min_max: numeric types only
  minlen_maxlen_len: string and slice
  required: see rule:check-required-semantics
  format_shortcuts: string primarily; skip empty unless required
  enum: comparable scalar or string
deferred_optional:
  - gt
  - gte
  - lt
  - lte
  - uri
  - url
  - format= generic sugar
excluded_v1:
  - eq / ne
  - cross-field
  - dive / element validation beyond outer length
  - unique
  - alpha / alphanum / contains family
  - file size / MIME (separate File rules later)
openapi_map:
  required: required / parameter required
  min: minimum
  max: maximum
  minlen: minLength or minItems
  maxlen: maxLength or maxItems
  len: minLength+maxLength or minItems+maxItems
  enum: enum
  pattern: pattern
  email: format email
  uuid: format uuid
  date: format date
  time: format time
  datetime: format date-time
  default: default
related:
  - concept:check-validation
  - rule:check-tag-syntax
  - rule:check-required-semantics
  - rule:check-format-validators
  - rule:openapi-validation-metadata
```
