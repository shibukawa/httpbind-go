---
id: rule:check-tag-syntax
type: rule
title: Check Tag DSL Syntax
---
check tag is a compact CSV-like DSL of rule tokens; enum values use pipe separators; pattern is trailing-only in v1.

```yaml
form: 'check:"rule,rule=value,..."'
token_kinds:
  - bare: required, email, uuid, date, time, datetime
  - key_value: min, max, minlen, maxlen, len, default, enum, pattern
separators:
  rules: ","
  enum_values: "|"
pattern_policy:
  v1: pattern= must be last token in the tag
  reason: commas inside regex break CSV split
  alternatives_deferred:
    - semicolon rule separators
    - quoted pattern values
enum_example: 'check:"required,enum=asc|desc|name"'
pattern_example: 'check:"required,pattern=^[A-Z]{3}$"'
not_compatible_with: go-playground/validator full dialect
parser: codegen only; never interpret tags at runtime
related:
  - concept:check-validation
  - rule:check-v1-rule-set
  - decision:check-tag-validation
  - decision:reflection-free
```
