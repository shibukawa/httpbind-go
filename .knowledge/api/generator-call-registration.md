---
id: api:generator-call-registration
type: api
title: Generator Call Registration API
---
Framework generator commands build an explicit local registry of data:generator-call-pattern values before analysis.

```yaml
public_shape:
  - "func NewCallRegistry() *CallRegistry"
  - "func (r *CallRegistry) Register(patterns ...CallPattern) error"
  - "func (r *CallRegistry) Options(base Options) (Options, error)"
construction:
  - operation-specific constructors define required semantic roles
  - ConfigBindCall maps config type plus prefix; ConfigSubCommandCall maps config type plus name and help
  - role source helpers select generic argument, value argument, argument type, or constant
  - SymbolPattern and MethodPattern constructors retain rule:go-types-symbol-identity
behavior:
  - registration is additive within one framework command construction
  - exact duplicates normalize to one entry
  - conflicting target semantics return an error
  - Options returns an immutable snapshot safe for concurrent package analysis
  - no package init, process-global registry, or runtime reflection
example: |
  calls := generator.NewCallRegistry()
  err := calls.Register(generator.ConfigBindCall(
      generator.Function("example.com/framework", "RegisterConfig"),
      generator.GenericType("config", 0),
      generator.Argument("prefix", 1),
  ))
  if err != nil { return err }
  options, err := calls.Options(generator.DefaultOptions())
related:
  - data:generator-call-pattern
  - data:generator-options
  - requirement:framework-wrapper-discovery
  - api:generator-main
```
