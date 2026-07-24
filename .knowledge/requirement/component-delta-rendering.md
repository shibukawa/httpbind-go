---
id: requirement:component-delta-rendering
type: requirement
title: Component Delta Rendering
---
Transmit only changed partial-update boundaries after page navigation or direct boundary execution.

```yaml
source: concept:html-render-runtime-extensions
flow: flow:html-partial-update
request:
  navigation: target path, new search parameters, and prior data:component-update-manifest validators
  boundary: data:boundary-parameter-update-request
execution_modes:
  navigation:
    - authenticate and validate the request normally
    - execute the generated route, runtime layouts, and page with current request state
    - build the next document update-boundary graph and validators
  boundary:
    - validate requirement:boundary-parameter-updates capability and typed changes
    - execute only the target boundary subtree with reconstructed immutable inputs
    - build the next subtree boundary graph and validators
common_execution:
  - use requirement:component-output-cache when eligible; otherwise render boundary output before comparison
comparison:
  match: rule:component-instance-identity
  unchanged: same content validator; omit boundary HTML
  changed: return replacement template and new validator
  structural: return insert, remove, or move; safely replace an ancestor when a granular operation is unavailable
response: data:component-delta-response
security:
  - client validators are optimization hints and cannot bypass authorization, validation, or page execution
  - server-derived current request values are the source of truth
  - do not reflect raw request data into operation scripts
failure:
  incompatible_render_version: full render or explicit reload response
  invalid_manifest: ignore hints and compute a safe full or larger delta
  render_error: no partially valid manifest publication
compatibility: requirement:html-rendering-compatibility
acceptance:
  - changing search parameters can re-execute the page without sending unchanged boundary HTML
  - changing an allowed data parameter can rerender only its declared boundary subtree
  - changed, inserted, removed, and reordered instances converge to the server render
  - next manifest represents the DOM state after all operations apply
open_questions:
  - request and response media types
  - validator hash algorithm, keyed mode, and wire compression
  - maximum manifest size and compact encoding
  - combining navigation deltas with requirement:suspense-html-streaming completions
```
