---
id: rule:check-required-semantics
type: rule
title: Check Required and Zero-Value Semantics
---
required must account for Go zero values; path/header absence differs from body/query scalars.

```yaml
problem: Go cannot distinguish omitted vs explicit zero for non-pointer scalars
v1_policy:
  string: required means non-empty
  slice: required means non-empty length
  path_header: missing extraction is required violation
  numeric_bool:
    prefer: pointer types or presence tracking when true required is needed
    v1_safe_default: allow required only on string/slice, or document zero-value pitfalls
  body_json:
    omit vs zero: same without pointer; accept limitation
format_interaction:
  empty_optional_fields: skip email/uuid/date/time/datetime when value empty and not required
  with_required: empty fails required before or instead of format
pipeline_note: rule:check-codegen-pipeline runs validate before default so optional absent skips min/max then may receive out-of-range sentinel defaults
sentinel_example:
  tag: 'check:"min=1,default=-1"'
  absent: after pipeline value is -1 (undefined to app)
  present_minus_one: fails min during validate; default not applied
related:
  - concept:check-validation
  - rule:check-v1-rule-set
  - rule:check-codegen-pipeline
  - concept:request-binding
```
