# SQL Template V1

Profile: `review`

| ID | Type | Title |
| --- | --- | --- |
| `requirement:sql-template-v1` | `requirement` | SQL Template V1 |
| `data:sql-statement` | `data` | Generated SQL Statement |
| `decision:postgresql-first-template-sql` | `decision` | PostgreSQL-First Template SQL |
| `decision:sql-dialect-generation-time` | `decision` | SQL Dialect Fixed at Generation Time |
| `decision:template-declaration-kinds` | `decision` | Format-Specific Template Declarations |
| `rule:sql-placeholder-emission` | `rule` | Generated SQL Placeholder Emission |
| `rule:template-name-casing` | `rule` | Template Name and Keyword Casing |
| `concept:typed-template-language` | `concept` | Compact Typed Template Language |
| `requirement:sql-generated-api-layers` | `requirement` | Generated SQL API Layers |
| `requirement:sql-relation-composition` | `requirement` | Typed SQL Relation Composition |
| `requirement:template-v1-scope` | `requirement` | Template V1 Scope |

## requirement:sql-template-v1

Generate parameterized SQL with typed result contracts and safe structured dynamic clauses.

```yaml
source: concept:typed-template-language
outputs:
  sql.exec: no row result; expose affected count when supported
  sql.one<T>: exactly one row; reject zero or multiple
  sql.optional<T>: zero or one row; reject multiple
  sql.many<T>: zero or more rows
  sql.predicate: reusable predicate list
  sql.relation<T>: private typed subquery relation from requirement:sql-relation-composition
declaration: lowercase statement keyword with PascalCase name from decision:template-declaration-kinds
naming: rule:template-name-casing
values: ordinary inserted expressions follow rule:sql-placeholder-emission
statement: data:sql-statement
generated_api: requirement:sql-generated-api-layers
dialect:
  initial: decision:postgresql-first-template-sql
  selection: decision:sql-dialect-generation-time
structured_lists:
  where: AND children by default; explicit and/or groups; omit when empty for SELECT
  joins: conditional; cannot vary result shape
  set: manage commas; require an unconditional item or pre-execution empty check
  order_by: static branches or enums; manage commas and empty clause
  insert: paired field-value assignments; no bulk insert
  returning: static item shape
relation_composition: requirement:sql-relation-composition
result_validation:
  - validate column count, names or aliases, types, optionality, and join nullability where provable
  - keep declared public cardinality when analysis is inconclusive
  - enforce unproven one and optional cardinality at runtime
mutation_safety:
  - UPDATE and DELETE reject an empty dynamic WHERE
  - full-table mutation needs a future explicit opt-in
forbidden:
  - value interpolation into SQL text
  - manually authored bind-placeholder tokens in executable SQL text; only value expressions generate them
  - arbitrary dynamic identifiers, operators, keywords, or sort directions
  - runtime-conditional select or returning columns
  - general loops in SQL clauses
```

## data:sql-statement

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

## decision:postgresql-first-template-sql

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

## decision:sql-dialect-generation-time

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

## decision:template-declaration-kinds

Use explicit lowercase declaration keywords instead of calling every output declaration a component.

```yaml
source:
  - concept:typed-template-language
  - user design discussion 2026-07-20
declarations:
  component:
    format: HTML
    required_output: html
  statement:
    format: SQL
    required_output_prefix: sql
common:
  - optional export modifier controls generated public API visibility
  - PascalCase declaration name
  - typed parameters and explicit output type
  - private declaration when export is absent
semantics:
  - declaration keyword selects format grammar
  - output type selects result behavior, insertion contexts, and SQL cardinality
  - keyword and output type mismatch is a compile-time error
compiler_model:
  common: TemplateDecl
  format_nodes: HTMLComponentDecl and SQLStatementDecl
```

## rule:sql-placeholder-emission

Template value expressions create bound arguments; template authors never manage database placeholders.

```yaml
source: concept:typed-template-language
template_syntax: 'column = {value}'
output: data:sql-statement
rules:
  - only a parsed template value expression may create a bind placeholder
  - raw SQL input never creates or preserves a bind-placeholder node or token
  - append evaluated value to Args and emit its placeholder as one operation
  - preserve encounter order after runtime structural conditions are resolved
  - reject manual bind-parameter tokens outside SQL literals and comments
  - never interpolate an ordinary value into SQL text
  - expand requirement:sql-relation-composition AST calls before numbering placeholders
parser:
  expression: lower to one generated placeholder plus one Args append
  sql_text: reject dialect-recognized bind-placeholder tokens
  literals_and_comments: placeholder-like characters remain inert source text
initial_styles:
  dollar_numbered: '$1, $2, ...'
  question: '?, ?, ...'
configuration:
  phase: code generation
  relation: separate option from SQL dialect; default follows decision:postgresql-first-template-sql
dynamic_sql: generated runtime appender owns numbering when optional clauses change argument count
```

## rule:template-name-casing

Case and word form identify language, user, HTML, and SQL symbol classes and are validated at compile time.

```yaml
source:
  - concept:typed-template-language
  - user design discussion 2026-07-20
classes:
  sql_keywords:
    form: UPPERCASE
    examples: [SELECT, FROM, LEFT JOIN, WHERE, IS NULL, TRUE, FALSE, NULL]
  dsl_keywords:
    form: lowercase
    examples: [export, component, statement, if, else, for, where, subquery, predicates]
  user_symbols:
    form: PascalCase
    includes: [types, enums, enum members, components, statements, external functions]
    examples: [UserRow, UserStatus, Active, UserCard, FindUser, NormalizeName]
  dsl_values:
    form: lowerCamelCase
    includes: [parameters, fields, local and loop variables]
    examples: [tenantID, minimumAge, displayName]
  sql_identifiers:
    form: lower_snake_case
    includes: [schemas, tables, columns, aliases]
    examples: [public, user_accounts, created_at, active_users]
  html_builtin_names:
    form: lowercase or kebab-case
    includes: [elements, attributes]
    examples: [div, aria-label, data-user-id]
  sql_builtin_names:
    form: lowercase unless classified as a dialect keyword
    includes: [functions, type names]
    examples: [count, coalesce, lower, integer]
  builtin_output_types:
    form: lowercase
    examples: [html, sql.exec, sql.many, sql.relation]
diagnostics:
  - recognize a SQL keyword written with wrong case and report required uppercase spelling
  - do not reinterpret wrong-case SQL keywords as identifiers
  - do not silently normalize user-defined symbols or format identifiers
  - user symbol resolution is case-sensitive and requires exact spelling
postgresql_v1:
  identifiers: lowercase unquoted only
  quoted_mixed_case_identifiers: deferred
```

## concept:typed-template-language

Small statically typed DSL for HTML composition and parameterized SQL generation. It is not a general-purpose language.

```yaml
evidence:
  source: user-supplied Compact Typed Template Language Specification
  received: 2026-07-20
review_gate: requirements remain proposed until user approval
outputs:
  - html
  - sql.exec
  - sql.one<T>
  - sql.optional<T>
  - sql.many<T>
  - sql.predicate
  - sql.relation<T>
principles:
  - parse output structure instead of interpolated raw strings
  - output type selects body parser, insertion rules, generated API, and SQL cardinality
  - keep exported component signatures explicit and stable
  - share typed declarations, expressions, and structural control across output formats
requirements:
  - requirement:template-language-core
  - requirement:sql-relation-composition
  - requirement:html-template-v1
  - requirement:explicit-output-control
  - requirement:sql-template-v1
  - requirement:sql-generated-api-layers
  - requirement:template-code-generation
  - requirement:template-v1-scope
boundary: decision:template-package-boundaries
safety: rule:template-context-safety
naming: rule:template-name-casing
declarations: decision:template-declaration-kinds
sql_dialect: decision:sql-dialect-generation-time
```

## requirement:sql-generated-api-layers

Generate a reusable statement builder and a database/sql execution wrapper for every exported executable SQL component.

```yaml
source: concept:typed-template-language
low_level:
  name: Build<Component>
  inputs: typed component parameters
  output: data:sql-statement plus error
  behavior: build SQL and Args without database access
high_level:
  name: <Component>
  inputs: context.Context, minimal executor interface, typed component parameters
  behavior: call low-level builder, execute, scan, and enforce declared result contract
executor_interfaces:
  sql.exec: ExecContext-compatible; accepts sql.DB, sql.Conn, and sql.Tx
  row_outputs: QueryContext-compatible; accepts sql.DB, sql.Conn, and sql.Tx
execution:
  sql.exec: ExecContext; return affected-row-capable result
  sql.one<T>: QueryContext; reject zero or multiple rows
  sql.optional<T>: QueryContext; accept zero or one and reject multiple rows
  sql.many<T>: QueryContext; scan all rows
query_row_rule: QueryRowContext is insufficient for multiple-row detection; use only when at-most-one is statically proven and the contract remains enforced
benefits:
  - low-level deterministic tests without a database
  - SQL logging, middleware, and custom execution
  - one generated public convenience API for normal database/sql use
```

## requirement:sql-relation-composition

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

## requirement:template-v1-scope

Keep the first implementation limited to the minimum language needed for safe HTML and SQL output composition.

```yaml
source: concept:typed-template-language
included:
  - requirement:template-language-core
  - requirement:sql-relation-composition
  - requirement:html-template-v1
  - requirement:explicit-output-control
  - requirement:sql-template-v1
  - requirement:sql-generated-api-layers
  - requirement:template-code-generation
deferred:
  - immutable let bindings
  - explicit enum underlying values and field mapping annotations
  - anonymous SQL row types if named rows suffice for the first milestone
  - typed SQL identifier abstraction and affected-row annotations
  - bulk insert and repeated SQL fragment syntax
excluded:
  - general user-defined value functions, lambdas, and pipelines
  - map, filter, reduce, mutable variables, generics, pattern matching, and macros
  - dynamic HTML names and arbitrary attribute spreads
  - block control inside HTML attribute values or attribute lists
  - arbitrary SQL identifier interpolation, general SQL loops, and dynamic result columns
  - async language semantics and runtime interpretation
milestone_order:
  - declaration keywords, naming rules, types, expressions, and signatures
  - HTML structure, escaping, components, if, and for
  - explicit raw output and safe script JSON contexts
  - SQL parameters, result contracts, and static statements
  - private typed relation statements in FROM and JOIN
  - PostgreSQL lowering, generated statement builders, and execution wrappers
  - structured SQL lists and mutation guards
```

## Review Checklist

- [ ] Scope is correct.
- [ ] Missing references are resolved.
- [ ] Policies and permissions are explicit.
- [ ] Generated output is not written back as source.
