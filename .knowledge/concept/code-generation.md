---
id: concept:code-generation
type: concept
title: Generated Runtime Code
---
Generator emits bind and write functions, validation, OpenAPI schemas, and streaming metadata without runtime reflection.

```yaml
artifacts:
  - request binding functions
  - response write functions
  - stream write functions
  - validation from concept:check-validation
  - OpenAPI schemas
  - streaming metadata
function_examples:
  - "func bindCreateUserRequest(r *http.Request) (CreateUserRequest, error)"
  - "func validateCreateUserRequest(v *CreateUserRequest) error"
  - "func writeCreateUserResponse(w http.ResponseWriter, r *http.Request, resp CreateUserResponse) error"
  - "func writeChatEventStream(w http.ResponseWriter, r *http.Request, stream httpbinder.Stream[ChatEvent]) error"
public_wrappers:
  - api:bind
  - api:write
  - api:write-error
discovery:
  - concept:handler-discovery
  - flow:handler-parse
  - rule:request-model-discovery
  - rule:response-model-discovery
  - rule:error-response-discovery
runtime: no reflection
related:
  - flow:code-generation
  - decision:reflection-free
  - concept:request-binding
  - concept:response-binding
  - concept:openapi-generation
  - concept:streaming
  - concept:stdlib-wrapper-unwrap
  - concept:check-validation
  - rule:check-codegen-pipeline
```
