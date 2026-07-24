---
id: rule:client-mutable-component-parameters
type: rule
title: Client-Mutable Component Parameters
---
Permit direct boundary updates only for component inputs explicitly declared mutable by the client update API.

```yaml
source: requirement:boundary-parameter-updates
declaration:
  boundary: component must enable requirement:partial-update-boundaries
  parameter: each client-mutable name is explicitly allowlisted with its template type
invocation:
  api: api:client-component-update
  params: object keys must match the generated mutable-parameter allowlist
  value: client encoder sends API values; generated Go decoder validates the declared type
constraints:
  - DOM attributes and arbitrary object keys cannot make an undeclared parameter mutable
  - reject unknown names, duplicate conflicting bindings, invalid values, and configured size-limit violations
  - immutable parent inputs, request identity, authorization state, and server-only values cannot be overwritten
  - direct rerender must not perform mutations or depend on exactly-once side effects
types:
  initial: string, bool, numeric primitives, enum, and optional forms supported by deterministic HTML controls
  deferred: records, file data, and arbitrary client objects unless an explicit codec exists
open_questions:
  - mutable declaration syntax
  - optional future declarative event binding under a collision-resistant reserved namespace
  - API debounce and multi-parameter batching options
  - validation error rendering within the boundary
```
