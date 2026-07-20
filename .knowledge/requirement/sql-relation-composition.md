---
id: requirement:sql-relation-composition
type: requirement
title: Typed SQL Relation Composition
---
Allow a private SQL statement to expose a typed row relation reusable as a structurally embedded subquery.

```yaml
source:
  - concept:typed-template-language
  - user design discussion 2026-07-20
declaration:
  keyword: statement
  output: sql.relation<T>
  visibility: private in v1; no generated execution API
invocation:
  from: FROM subquery RelationName(args) AS alias
  join: JOIN subquery RelationName(args) AS alias
  alias: required lower_snake_case identifier
typing:
  - T is a named static row type
  - selected columns must match T
  - outer references through alias are checked against T
  - runtime-conditional result columns are forbidden
composition:
  - inline referenced relation AST into the outer typed SQL AST
  - resolve all explicit relation arguments in caller scope
  - perform dialect lowering after relation expansion
  - emit placeholders and Args once across the expanded statement via rule:sql-placeholder-emission
constraints:
  - no implicit correlated reference to outer aliases
  - no recursive statement calls
  - no direct SQL string or data:sql-statement concatenation
deferred:
  - sql.scalar<T>
  - CTE declaration and reuse
  - correlated and LATERAL subqueries
  - recursive CTEs
  - cross-module public relation fragments
```
