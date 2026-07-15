---
id: rule:wrapper-strip-prefix
type: rule
title: StripPrefix Unwrap
---
Unwrap http.StripPrefix and record stripped prefix; OpenAPI path stays the mux pattern.

```yaml
wrapper: "http.StripPrefix(prefix, h)"
example: |
  mux.Handle(
      "POST /api/",
      http.StripPrefix(
          "/api",
          http.HandlerFunc(apiHandler),
      ),
  )
parse:
  - unwrap h
  - analyze inner handler
  - record stripped prefix
runtime_note: StripPrefix modifies request path passed to inner handler
openapi:
  path_source: mux registration pattern
  does_not: automatically create a new OpenAPI path
  metadata: stripped_prefix
related:
  - concept:stdlib-wrapper-unwrap
  - concept:route-discovery
  - concept:openapi-generation
```
