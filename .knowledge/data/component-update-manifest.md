---
id: data:component-update-manifest
type: data
title: Component Update Manifest
---
Compact client-held state for comparing updateable component instances between renders.

```yaml
source: concept:html-render-runtime-extensions
document:
  route_id: stable generated route or page identity
  render_version: generated-code and update-protocol version
  document_validator: optional whole-render validator
instances:
  - instance_id: rule:component-instance-identity
    parent_id: optional enclosing update boundary
    component_id: stable generated declaration identity
    revision: monotonic boundary state version
    input_validator: opaque digest of component version and canonical typed inputs
    content_validator: opaque digest of canonical rendered boundary HTML
    position: optional structural ordering token
    direct_update:
      handle: stable public boundary reference accepted by api:client-component-update
      continuation: opaque authenticated token or server-state reference
      mutable_parameters: rule:client-mutable-component-parameters descriptors
transport:
  initial_html: encoded near boundary markers or in one manifest payload
  update_request: client returns prior render version and instance validators
privacy:
  - omit raw component arguments by default
  - validators may be keyed or opaque when plain hashes expose sensitive low-entropy values
  - server reconstructs arguments from the new request and generated render plan
  - continuation never replaces current authentication or authorization checks
canonicalization:
  - exclude compression and request-unique transport markers from content hashing
  - include component generated-code version in input validation
```
