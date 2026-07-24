---
id: data:boundary-parameter-update-request
type: data
title: Boundary Parameter Update Request
---
Typed request to rerender one declared component boundary without reexecuting its page or ancestors.

```yaml
source: requirement:boundary-parameter-updates
fields:
  render_version: generated code and protocol compatibility
  instance_id: rule:component-instance-identity
  base_revision: last applied boundary revision
  continuation: opaque capability from data:component-update-manifest
  changes:
    parameter_name: wire value allowed by rule:client-mutable-component-parameters
  known_validators: target subtree entries from data:component-update-manifest
semantics:
  - changes are patches over reconstructed current boundary inputs
  - omission preserves an existing input; explicit null follows the declared optional type
  - request does not carry arbitrary component identity or immutable arguments
safety:
  - apply normal authentication, authorization, origin or CSRF policy, request limits, and typed validation
  - enforce policy:html-update-csrf-protection before decoding continuation or rendering
  - capability binds route, component, instance, immutable inputs, generated version, and expiry or server session
  - stale base revision is rejected, rebased by policy, or answered with authoritative replacement
```
