# httpbind-go (`httpbinder`)

[日本語](README.ja.md)

Reflection-free, code-generation-first library that bridges Go types and HTTP APIs.

Define request/response structs once. The generator emits type-specific binders and writers, so the same model covers **JSON, form, multipart, and query** (plus path / header / cookie via tags). Responses adapt to the client **`Accept`** (and streaming negotiation where used). From the same analysis it also **generates OpenAPI 3.1**, kept in sync with binders and writers. Route registration is discovered by **static analysis of real `net/http` styles** (`HandleFunc`, `Handle`, method values, wrappers, and so on)—not by a separate DSL.

```go
type CreateUserRequest struct {
	// input = query + payload (JSON / form / multipart). Tag may be omitted.
	Name  string `input:"name"`  // same as untagged: Name string
	Email string `input:"email"` // same as untagged: Email string
	OrgID string `path:"org_id"`
	Token string `header:"Authorization"`
}

type CreateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	OrgID string `json:"org_id"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	input, err := httpbinder.Bind[CreateUserRequest](r)
	if err != nil {
		httpbinder.WriteError(w, r, err)
		return
	}
	// Name/Email: query and/or JSON/form/multipart body (input).
	// OrgID from path, Token from Authorization header.
	out := CreateUserResponse{
		ID:    "u_1",
		Name:  input.Name,
		Email: input.Email,
		OrgID: input.OrgID,
	}
	_ = httpbinder.Write[CreateUserResponse](w, r, out)
}
```

Run the generator on the package (binders + OpenAPI embed):

```bash
go run ./cmd/httpbinder-gen -dir . -openapi
```

### Struct tag reference

Wire name defaults to the lower-camel field name when a tag value is omitted (e.g. untagged `Name` → `"name"`).

| Tag | Source | Notes |
|-----|--------|--------|
| *(none)* or `input:"name"` | **query + payload** | Default. Payload covers JSON, `application/x-www-form-urlencoded`, and `multipart/form-data`. Tag is optional when the field is plain user input. |
| `query:"page"` | query only | Not read from the body. |
| `payload:"name"` | body only | JSON / form / multipart by `Content-Type`. Not read from the query string. |
| `payload:"image"` on `httpbinder.File` | multipart file part | Binds filename, content type, size, and bytes from the named part. Payload-only (not query). Multipart bodies are capped at **1 MiB** by default; override with `httpbinder.SetMaxMultipartBodyBytes`. |
| `path:"org_id"` | path parameter | Matches `{org_id}` (or equivalent) in the route pattern. |
| `header:"Authorization"` | request header | Header name is the tag value. |
| `cookie:"session"` | cookie | Cookie name is the tag value. |

**`input` vs `payload` vs `query`**

- Prefer **`input`** (or no tag) for normal fields that may arrive as query *or* body.
- Use **`query`** / **`payload`** only when you must restrict the origin (e.g. search filters in the query string, body-only JSON fields).
- `payload` is not the same as `input`: it does **not** accept query parameters.

Example that mixes restrictions:

```go
type SearchRequest struct {
	Keyword string `query:"keyword"`   // query only
	Page    int    `query:"page"`
	Filter  string `payload:"filter"`  // body only (JSON/form/multipart)
}
```

Response structs commonly use standard `json:"..."` names for encoding; request binding still uses the source tags above.

### Streaming (ideal API)

```go
stream, err := httpbinder.NewStream[ChatEvent](w, r)
if err != nil {
    httpbinder.WriteError(w, r, err)
    return
}
defer stream.Close()

_ = stream.Write(ChatEvent{Type: "delta", Delta: "hi"})
_ = stream.Write(ChatEvent{Type: "done"})
```

- **`Write` can be called many times** (incremental events).
- Format is chosen once in `NewStream` from `?stream=`, `Accept`, `User-Agent`, then default **NDJSON**.
- Formats:
  - **SSE** — `text/event-stream`
  - **NDJSON / JSONL** — `application/x-ndjson` (one object per line; *not* a JSON array)
  - **JSON array** — `application/json` as `[obj1,obj2,...]` (`Close` writes the trailing `]`)
- Do **not** use removed helpers `WriteNDJSON` / `WriteSSE`.

## Packages

| Path | Role |
|------|------|
| `.` (`package httpbinder`) | Runtime: Bind / Write / WriteError / NewStream / OpenAPI serve / SwaggerUI |
| `generator/` | Field-plan binders/writers + OpenAPI 3.1 embed generation |
| `parser/` | Route/handler discovery (`Bind`, `Write`, `NewStream`, errors) |
| `cmd/httpbinder-gen` | CLI: binders + OpenAPI from a package dir |
| `examples/demo` | End-to-end sample app |
| `internal/*` | Test fixtures |
| `testdata/cmd/*` | Dev-only helpers (not for distribution; under `testdata` so `go get` / `./...` skip them) |

```bash
go run ./cmd/httpbinder-gen -dir ./path/to/package
```

## Demo

```bash
go generate ./examples/demo
go run ./examples/demo
# http://localhost:8080/       index + browser stream demo
# http://localhost:8080/docs/  Swagger UI
# http://localhost:8080/chat   NewStream (SSE / NDJSON / JSON array auto)
```

See [`examples/demo/README.md`](examples/demo/README.md) for full curl recipes.

## TinyGo

TinyGo is a design goal for the reflection-free binder path. See notes below for toolchain limits.

Verified with **TinyGo 0.40.1** (Go **1.19–1.25**). System Go 1.26 is rejected by TinyGo 0.40.

```bash
./scripts/tinygo-check.sh
```

### Runtime notes relevant to TinyGo

- `AsHTTPError` avoids `errors.As` (unimplemented `AssignableTo` on some TinyGo builds).
- `WriteError` hand-builds problem JSON (avoids fragile nested `encoding/json` + RawMessage interactions).
- Registry uses `reflect.Type` only as a **type identity key**, not for field walking.
- Generated bind/write code does not import `reflect`.

### Known limitations

| Topic | Limitation |
|-------|------------|
| Toolchain | TinyGo 0.40 needs Go ≤ 1.25 (`GOTOOLCHAIN=go1.25.4`) |
| Streaming | Prefer host `go test` for `NewStream`; not fully TinyGo-matrixed |
| ServeMux | Prefer testing handlers with `ServeHTTP` + `SetPathValue` under TinyGo |
| Multipart `File` | Supported via `httpbinder.File` (`payload`); size/MIME `check` rules deferred. Body cap defaults to **1 MiB** (`SetMaxMultipartBodyBytes`) |
| Generator | Host-side only (`go run` / `go test`) |

## License

Licensed under the [Apache License, Version 2.0](LICENSE).
