---
id: api:config-scaffold-output
type: api
title: Config Scaffold Output API
---
Public configbind functions render scaffold fields from all registered configbind Definition values without owning a CLI or file path.

```yaml
signatures:
  - func ScaffoldTOML() (string, error)
  - func ScaffoldEnv() (string, error)
  - func WriteScaffoldTOML(w io.Writer) error
  - func WriteScaffoldEnv(w io.Writer) error
errors:
  - invalid definition metadata
  - duplicate config key
  - duplicate environment name
requirement: requirement:scaffold-generation
```
