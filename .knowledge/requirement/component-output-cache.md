---
id: requirement:component-output-cache
type: requirement
title: Component Output Cache
---
Reuse explicitly cache-enabled component output for equivalent typed inputs until its TTL expires.

```yaml
source: concept:html-render-runtime-extensions
declaration:
  activation: explicit component cache flag
  options: TTL required; additional policy syntax unresolved
identity:
  key_parts:
    - stable component identity and generated-code version
    - canonical type-aware encoding of every declared input parameter
  equality: different types or parameter boundaries cannot alias
value: validated rendered HTML bytes for one component invocation
behavior:
  hit: write cached bytes through the caller's current response stream
  miss: render once, publish only after successful complete rendering, then reuse until expiry
  expiry: expired entries behave as misses
partial_update:
  boundary: independent from requirement:partial-update-boundaries; cache and client update flags may be enabled separately
  reuse: cached content validator may satisfy requirement:component-delta-rendering without rerendering
safety:
  - never cache errors or partial output
  - preserve rule:template-context-safety at insertion
  - cache only output whose complete dependency set is represented by the key
  - request, authorization, locale, and header-derived variation must be explicit inputs or disable caching
compatibility: requirement:html-rendering-compatibility
acceptance:
  - identical inputs within TTL can avoid component logic and rendering
  - changed input, component version, or expired TTL cannot reuse the old value
open_questions:
  - process-local versus pluggable shared store
  - eviction and maximum entry or byte limits
  - concurrent miss coalescing and cancellation ownership
  - stale-while-revalidate and explicit invalidation
  - cache interaction with requirement:suspense-html-streaming
```
