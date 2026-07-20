---
id: decision:sql-dialect-generation-time
type: decision
title: SQL Dialect Fixed at Generation Time
---
Select SQL dialect and placeholder style when running the code generator, never when executing generated application APIs.

```yaml
source:
  - concept:typed-template-language
  - user design discussion 2026-07-20
generator_options:
  dialect: postgresql or future sqlite
  placeholder_style: rule:sql-placeholder-emission
pipeline:
  - parse to dialect-neutral typed SQL IR
  - validate selected dialect capabilities
  - lower dialect-specific types and syntax
  - bake placeholder appender into generated code
  - emit requirement:sql-generated-api-layers
runtime:
  receives:
    - component parameters
    - runtime structural condition values
    - database executor for high-level API
  excludes:
    - dialect argument
    - placeholder-style argument
    - driver-based dialect detection
multi_dialect: generate separate packages or artifacts for each dialect
benefits:
  - deterministic SQL and golden tests
  - generation-time unsupported-feature diagnostics
  - no per-query dialect branching
  - stable generated public APIs
```
