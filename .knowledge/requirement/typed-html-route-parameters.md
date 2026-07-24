---
id: requirement:typed-html-route-parameters
type: requirement
title: Typed HTML Route Parameters
---
Bind bracketed route segments to statically declared page and layout parameters.

```yaml
source: concept:filesystem-html-routing
discovery: decision:html-route-file-conventions
example: routes/users/[id]/index.tb.html maps one URL segment to id
typing:
  name: bracketed directory supplies the parameter name
  type: page or layout declaration supplies a supported scalar type; string is the simplest default
  binding: generated decoder validates URL segment before rendering
scope:
  page: may accept every dynamic segment from route root through page directory
  layout: may accept only dynamic segments from route root through its own directory
  parent_layout: cannot depend on deeper child segment parameters
search_parameters:
  page: may declare a distinct typed search-parameter input
  layout: excluded by default so search changes do not invalidate ancestor wrappers
validation:
  - declaration name must match an in-scope dynamic segment
  - duplicate parameter names in one route are generation errors
  - missing required declaration inputs and incompatible types are generation errors
  - invalid request value returns configured bad-request or not-found behavior before rendering
open_questions:
  - supported non-string scalar decoders
  - declaration form for optional page search parameters
  - invalid path value status policy
```
