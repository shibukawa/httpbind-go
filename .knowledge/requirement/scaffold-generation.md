---
id: requirement:scaffold-generation
type: requirement
title: Config Scaffold Generation
---
Codegen exposes Bind-based TOML and .env scaffold text as generated Go constants; applications own any printing command.

```yaml
priority: must
intent: bootstrap shared config files from Bind structs only
delivery: exported constants in the generated Bind Go file
mechanism:
  - codegen embeds plain text scaffold bodies from struct metadata
  - ConfigbindScaffoldTOML contains TOML text
  - ConfigbindScaffoldEnv contains .env text
  - application code may print or write the constants through its chosen CLI
inputs:
  - api:configbind-bind registrations only
  - decision:struct-field-tags for default, help, enum
  - data:cli-flag-def help text for comments
  - decision:prefix-table-binding
  - decision:toml-shape-constraints
excluded_inputs:
  - api:configbind-subcommand types and fields
outputs:
  - generated TOML string constant with prefix tables, dotted nested keys, and primitive arrays
  - generated .env string constant using runtime environment naming and overrides
  - comments derived from help tags next to keys
  - example values derived from default tags when present
  - optional allowed-value notes from enum tags
constraints:
  - codegen performs no runtime file write
  - codegen adds no application CLI command or subcommand
  - do not emit inline tables, arrays of tables, or quoted keys
  - nested structs become nested tables or dotted bare keys
  - do not include subcommand options or positionals
  - opt CLI renames do not change TOML key names in the scaffold
related:
  - flow:configbind-codegen
  - concept:scaffold-templates
  - requirement:struct-field-metadata
  - requirement:cli-option-codegen
  - requirement:cli-subcommands
  - decision:struct-field-tags
  - decision:cli-flag-naming
  - data:cli-flag-def
  - system:configbind
acceptance:
  - generated package exports ConfigbindScaffoldTOML and ConfigbindScaffoldEnv
  - application-owned code can print either constant unchanged
  - scaffold TOML contains [prefix] tables for each Bind
  - scaffold TOML lines include help as comments when help tag present
  - default values appear as example values when default tag present
  - scaffold env uses runtime names including opt and env overrides and omits env:"-" fields
  - subcommand-only fields never appear in scaffolds
```
