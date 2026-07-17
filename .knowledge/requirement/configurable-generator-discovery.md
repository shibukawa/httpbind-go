---
id: requirement:configurable-generator-discovery
type: requirement
title: Configurable Generator Discovery
---
Applications may construct a generator/parser with compatible symbol identities instead of the built-in net/http and httpbinder set.

```yaml
status: required
configuration:
  call_identity: package path + function name
  method_identity: package path + method name + receiver package path + receiver type
  special_type_identity: package path + type name + generator role
defaults:
  calls: Bind, Write, WriteStatus, DecodeJSON, EncodeJSON, NewStream
  routes: net/http package registrars and net/http.ServeMux methods
  special_types: httpbinder.File
behavior:
  - zero configuration uses defaults
  - explicit configuration can replace or extend defaults
  - aliases resolve through rule:go-types-symbol-identity
  - same-named unconfigured symbols never match
public_surface:
  - reusable configured generator object
  - configured package analyzer
  - configured route parser
host_only: true
related:
  - requirement:strict-symbol-identity
  - rule:go-types-symbol-identity
  - concept:code-generation
  - flow:code-generation
```
