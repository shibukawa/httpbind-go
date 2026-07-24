---
id: decision:html-route-file-conventions
type: decision
title: HTML Route File Conventions
---
Discover pages and layouts from reserved filenames and bracketed path-segment directories below an opt-in route root.

```yaml
source: concept:filesystem-html-routing
review_gate: proposed convention requires user approval
activation: generator configuration explicitly selects one route root
preferred_files:
  page: index.tb.html
  layout: layout.tb.html
  document: document.tb.html only at configured route root
logical_aliases: index.html, layout.html, and document.html describe roles; accepting them as physical files remains open
segments:
  static: directory name becomes a literal URL segment
  dynamic: '[id]' becomes one named path parameter through requirement:typed-html-route-parameters
  catch_all: deferred
precedence:
  - literal route outranks a dynamic sibling
  - duplicate normalized route patterns are generation errors
  - at most one dynamic sibling exists at a directory level unless patterns are provably distinct
tree:
  page: index file owns the exact directory route
  layouts: ancestor layout files ordered from route root to page directory
  document: optional root-only shell from decision:html-document-shell; never owns the root URL
constraints:
  - route tree templates belong to one configured generation unit; child directories are not implicitly separate Go packages
  - Go output and external binding follow decision:html-route-go-package-model
  - reserved behavior does not apply to existing flat template discovery
  - generated identifiers include normalized relative path to prevent filename collisions
open_questions:
  - physical extension aliases and route-root configuration syntax
  - route groups that do not contribute URL segments
  - optional and catch-all segment notation
```
