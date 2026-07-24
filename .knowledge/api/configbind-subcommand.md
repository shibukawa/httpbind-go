---
id: api:configbind-subcommand
type: api
title: configbind.SubCommand
---
Generic SubCommand registers a CLI-only subcommand option struct and returns *T when selected, nil otherwise.

```yaml
signature_sketch: 'func SubCommand[T any](name string, help string) *T'
return_type: '*T'
nil_semantics:
  - non-nil *T when CLI selected this subcommand and parse succeeded
  - nil when not selected
  - nil when selected arguments fail parsing; Load returns the stored UsageError
behavior:
  - generated code registers SubCommandDefinition before main
  - SubCommand selects from process argv and parses CLI flags and positionals on T only
  - Load validates the same branch and returns generated usage for help or errors
  - does not participate in TOML or env layers
  - parses positional fields tagged arg required|optional|*
  - options may appear before or after positionals
testing:
  - LoadOptions.Args must match os.Args[1:] because nil/non-nil selection occurs when SubCommand returns
example:
  go: |
    migrate := configbind.SubCommand[MigrateOpt]("migrate", "run migrations")
    // app migrate ./db --dry-run  -> migrate != nil
    // other subcommand or none    -> migrate == nil
depends_on:
  - requirement:cli-subcommands
  - decision:struct-field-tags
  - decision:configbind-supported-types
related:
  - api:configbind-bind
  - concept:subcommand-binding
  - flow:config-load
  - requirement:cli-option-codegen
  - system:configbind
```
