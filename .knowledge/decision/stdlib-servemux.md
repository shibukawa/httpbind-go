---
id: decision:stdlib-servemux
type: decision
title: Stdlib ServeMux Routing
---
Route registration uses Go 1.22+ net/http ServeMux; no framework-specific router is required.

```yaml
status: accepted
router: net/http ServeMux
go_version_min: "1.22"
example: |
  mux := http.NewServeMux()
  mux.HandleFunc(
      "POST /orgs/{org_id}/users",
      CreateUserHandler,
  )
path_params:
  - from: route patterns such as {org_id}
    to: path-tagged request fields
forbidden_requirement:
  - framework-specific router
related:
  - concept:net-http-handler
  - term:http-metadata
  - concept:handler-discovery
  - concept:route-discovery
  - rule:unsupported-route-patterns
```
