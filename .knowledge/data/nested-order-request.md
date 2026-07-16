---
id: data:nested-order-request
type: data
title: NestedOrderRequest Example
---
Example nested request using struct, slice, and map fields under rule:nested-request-binding.

```yaml
type: |
  type NestedOrderRequest struct {
      Customer NestedCustomer            `payload:"customer"`
      Items    []NestedLineItem          `payload:"items"`
      Labels   map[string]string         `payload:"labels"`
  }
  type NestedCustomer struct {
      ID   string `payload:"id"`
      Name string `payload:"name"`
  }
  type NestedLineItem struct {
      SKU string `payload:"sku"`
      Qty int    `payload:"qty"`
  }
json_example: |
  {
    "customer": { "id": "c1", "name": "Ada" },
    "items": [
      { "sku": "A-1", "qty": 2 },
      { "sku": "B-9", "qty": 1 }
    ],
    "labels": { "channel": "web", "priority": "high" }
  }
binding:
  Customer: nested object via rule:nested-request-binding
  Items: slice of structs
  Labels: map[string]string
openapi: rule:openapi-nested-schemas
related:
  - rule:nested-request-binding
  - term:payload
  - concept:request-binding
  - concept:openapi-generation
```
