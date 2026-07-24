---
id: data:async-render-error
type: data
title: Async Render Error
---
Safe typed failure value exposed to an asynchronous boundary recover clause.

```yaml
source: requirement:async-external-functions
public_fields:
  code: stable application or generated classification
  message: optional presentation-safe text
  retryable: whether UI may offer a retry
  timeout: whether configured async deadline expired
server_only:
  - original Go error or panic value
  - stack and component call chain
  - request and boundary diagnostic context
normalization:
  returned_error: map through declared or default classifier
  panic: log raw value and expose generic internal code
  timeout: expose stable timeout code and retry policy
  cancellation: expected request or supersession cancellation produces no recover value
constraints:
  - raw error text is not public by default
  - error value is immutable and usable only inside decision:async-boundary-syntax recover subtree
  - failure results and recover HTML are not stored by requirement:component-output-cache
```
