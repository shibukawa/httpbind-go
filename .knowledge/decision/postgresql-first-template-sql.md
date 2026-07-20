---
id: decision:postgresql-first-template-sql
type: decision
title: PostgreSQL-First Template SQL
---
Use PostgreSQL as the first SQL semantic target while keeping the initial AST and feature subset portable.

```yaml
source:
  - concept:typed-template-language
  - user design discussion 2026-07-20
default:
  dialect: postgresql
  placeholder: dollar_numbered from rule:sql-placeholder-emission
rationale:
  - strict and rich database types align with static template types
  - schema and result validation can be stronger
  - PostgreSQL supports the planned returning and structured mutation workflows
portable_v1:
  - SELECT, INSERT, UPDATE, and DELETE
  - joins, where, order by, limit, and offset
  - basic returning
  - bound values and expanded IN placeholders
future_sqlite:
  priority: second dialect before broad PostgreSQL-only language features
  requires:
    - dynamic-affinity and STRICT-table schema handling
    - explicit date, time, datetime, decimal, and boolean storage mappings
    - placeholder expansion and parameter-limit checks
    - RETURNING capability restrictions
future_postgresql:
  optional_lowering:
    - array parameters and ANY
    - native JSON and JSONB
    - richer returning and PostgreSQL-specific types
constraint: dialect-specific syntax requires capability validation and must not silently change portable semantics
```
