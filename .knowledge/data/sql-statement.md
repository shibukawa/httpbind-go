---
id: data:sql-statement
type: data
title: Generated SQL Statement
---
Transport-neutral low-level result of a generated SQL component before database execution.

```yaml
source: concept:typed-template-language
go_shape:
  SQL: string
  Args: '[]any'
properties:
  - SQL contains only generator-owned bind placeholders
  - Args follow placeholder emission order
  - no database handle, rows, or dialect selection
construction_errors:
  - unsafe empty mutation WHERE
  - empty dynamic SET
  - other runtime-dependent structural validation failures
```
