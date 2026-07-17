---
id: requirement:strict-symbol-identity
type: requirement
title: Strict Symbol Identity for Static Analysis
---
Host-side generator analysis must identify net/http and tinybind symbols by resolved type identity (go/types), not bare selector names.

```yaml
status: implemented
problem:
  - route discovery accepts any receiver named Handle or HandleFunc
  - handler body may mis-detect Bind/Write from unrelated packages with the same names
  - error constructors fail recognition when tinybind is imported under an alias
approach: rule:go-types-symbol-identity
configuration:
  default: net/http and tinybind identities
  custom: requirement:configurable-generator-discovery
not: name-string matching as source of truth
host_only: true
note: uses go/types on host Go toolchain; generator remains host-side only
related:
  - rule:go-types-symbol-identity
  - concept:handler-discovery
  - concept:route-discovery
  - flow:handler-parse
  - flow:code-generation
  - concept:code-generation
```
