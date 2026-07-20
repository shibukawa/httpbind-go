---
id: decision:template-parser-delegation
type: decision
title: Delegated Template Parser Architecture
---
Keep format syntax out of the shared parser by recursively delegating declaration bodies to registered format parsers.

```yaml
source:
  - concept:typed-template-language
  - user architecture discussion 2026-07-20
shared_parser:
  owns:
    - modules, imports, types, enums, external functions, and declaration headers
    - expression and if/else/for header grammars
    - common control AST construction, source positions, scopes, and diagnostics
    - recursive orchestration between control bodies and the active format parser
  excludes:
    - HTML tags, attributes, raw-text rules, and escaping contexts
    - SQL tokens, statements, clauses, comments, literals, and placeholders
    - scanning arbitrary format text for embedded-template boundaries
format_parser:
  owns:
    - format tokenization and structure
    - discovery of expression and if/for boundaries in its current lexical context
    - selection of the insertion or structural context for each shared node
    - returning control terminators to the shared parser
  calls_shared:
    - parse an expression after discovering an expression boundary
    - parse an if/for header after discovering a control boundary
    - construct shared control nodes and coordinate branches
  recursion:
    - shared parser sends each then, else, and loop body back to the same format parser
    - format parser stops at a matching control terminator and returns control to the shared parser
registration:
  key: lowercase root declaration keyword
  examples:
    component: HTML format parser
    statement: SQL format parser
  constraint: shared parser contains no HTML- or SQL-specific branches or imports
test_dummy_parser:
  role: test one complete shared-parser integration without HTML or SQL knowledge
  behavior:
    - preserve non-template input as lossless raw text
    - discover expression and if/for boundaries in raw text
    - delegate shared grammar and recursively parse control bodies
    - attach raw:text context to embedded expressions and controls
ast_type_ids:
  format: '<language-code>:<node-type>'
  shared: [template:expression, template:if, template:for]
  dummy: [raw:text]
  html: [html:element, html:attribute, html:text, html:component]
  sql: [sql:statement, sql:clause, sql:parameter]
contexts:
  representation: opaque namespaced string owned by the format parser
  examples: [raw:text, html:child, html:attribute, html:script, html:style, sql:value]
  rule: every shared expression and control node records the context supplied by its format parser
source_positions:
  shape: 'pos: {line: int, col: int}'
  basis: one-based file-global Unicode character position
```
