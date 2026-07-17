---
id: rule:sql-group-key
type: rule
title: SQL Join Group Keys
---
The groupkey struct tag defines identity boundaries when joined rows become parent-child object trees.

```yaml
status: required
tag: 'groupkey:""'
scope: each root or slice-element struct participating in grouping
constraints:
  - exactly one scalar groupkey per grouped struct
  - groupkey field also maps from its db column
  - duplicate keys merge even when rows are not contiguous
  - root order and child order follow first appearance
  - scalar values use the first row for each object
null:
  root_key: error
  child_key: omit child for outer join
nested_slices: recursive grouping with the same rules
non_slice_structs: flatten scalar columns into the containing object
related:
  - api:scan-rows
  - rule:usage-directed-generation
  - decision:reflection-free
```
