---
id: rule:check-format-validators
type: rule
title: Check Format Validators
---
Format shortcuts use pragmatic fixed checks; date/time are ISO-only; email is intentionally non-strict.

```yaml
date:
  accept: "2006-01-02"
  go: time.DateOnly
  openapi_format: date
  reject: non-ISO layouts
time:
  accept: "15:04:05"
  go: time.TimeOnly
  openapi_format: time
datetime:
  tag_name: datetime
  accept: RFC3339
  optional_fallback: RFC3339Nano after RFC3339 fail
  reject: timezone-less local datetimes
  openapi_format: date-time
  naming: tag uses datetime; OpenAPI uses date-time
uuid:
  intent: valid UUID string
  openapi_format: uuid
email:
  strictness: pragmatic not RFC5322
  checks:
    - non-empty when required
    - exactly one '@'
    - non-empty local and domain
    - no whitespace
    - domain contains at least one '.'
  combine_with: maxlen=254 recommended
  openapi_format: email
  escape_hatch: user may add pattern for stricter rules
  avoid: net/mail.ParseAddress as sole check (accepts display-name forms)
pattern:
  engine: Go regexp RE2
  tinygo: regexp supported
  codegen:
    - compile at generation time; invalid pattern is codegen error
    - emit package-level compiled regexp vars
  limits: no backrefs or lookahead; document simple constraints preferred
  syntax: rule:check-tag-syntax
related:
  - concept:check-validation
  - rule:check-v1-rule-set
  - requirement:tinygo-wasm
  - rule:openapi-validation-metadata
```
