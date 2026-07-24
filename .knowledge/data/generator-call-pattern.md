---
id: data:generator-call-pattern
type: data
title: Generator Call Pattern
---
A call pattern maps one framework function or method identity to a generator semantic operation without inspecting wrapper implementation.

```yaml
CallPattern:
  target: SymbolPattern or MethodPattern from data:generator-options
  operation:
    - request_bind
    - response_write
    - response_write_status
    - stream_create
    - json_decode
    - json_encode
    - rows_scan
    - config_bind
    - config_subcommand
    - route_register
    - error_response
  type_roles: map of semantic role to TypeSource
  argument_roles: map of semantic role to ValueSource
TypeSource:
  generic_argument: zero-based generic parameter index
  argument_type: zero-based value argument index
ValueSource:
  argument: zero-based value argument index
  constant: statically configured scalar value
indexing:
  receiver: excluded from value argument indices
  variadic: its declared argument owns one index
behavior:
  - unmapped wrapper arguments are ignored
  - reordered operation arguments use explicit role indices
  - wrapper-added context, logging, and framework option arguments require no mapping
  - a hidden fixed status or name uses constant ValueSource
  - generic type inference is resolved through go/types Instances
validation:
  - target identity must resolve uniquely
  - every operation-required role has exactly one source
  - role indices must exist in the resolved wrapper signature
  - a source kind must satisfy the required role type
  - conflicting patterns for one target are rejected before package analysis
examples:
  config_wrapper:
    call: RegisterConfig[ServerConfig](ctx, "server")
    operation: config_bind
    type_roles: {config: {generic_argument: 0}}
    argument_roles: {prefix: {argument: 1}}
  subcommand_wrapper:
    call: Command[MigrateOptions](ctx, "migrate", "run migrations")
    operation: config_subcommand
    type_roles: {config: {generic_argument: 0}}
    argument_roles: {name: {argument: 1}, help: {argument: 2}}
  status_wrapper:
    call: WriteCreated[User](trace, w, r, user)
    operation: response_write_status
    type_roles: {response: {generic_argument: 0}}
    argument_roles:
      writer: {argument: 1}
      request: {argument: 2}
      value: {argument: 3}
      status: {constant: 201}
```
