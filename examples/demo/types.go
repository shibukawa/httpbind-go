package main

// CreateUserRequest demos default input, path, and header binding.
type CreateUserRequest struct {
	Name  string
	Email string
	OrgID string `path:"org_id"`
	Token string `header:"Authorization"`
}

// CreateUserResponse is returned after create.
type CreateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	OrgID string `json:"org_id"`
}

// SearchRequest demos query-only and payload-only fields.
type SearchRequest struct {
	Keyword string `query:"keyword"`
	Page    int    `query:"page"`
	Filter  string `payload:"filter"`
}

// SearchResponse is search output.
type SearchResponse struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Filter  string `json:"filter"`
	Hits    int    `json:"hits"`
}

// EchoRequest demos input from query or form/JSON body.
type EchoRequest struct {
	Message string
	N       int `query:"n"`
}

// EchoResponse echoes back.
type EchoResponse struct {
	Message string `json:"message"`
	N       int    `json:"n"`
}

// SessionRequest demos cookie binding.
type SessionRequest struct {
	Session string `cookie:"session"`
}

// SessionResponse shows cookie value.
type SessionResponse struct {
	Session string `json:"session"`
	OK      bool   `json:"ok"`
}

// UserGetRequest path-only lookup.
type UserGetRequest struct {
	ID string `path:"id"`
}

// UserGetResponse user payload.
type UserGetResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ChatRequest starts a demo stream.
type ChatRequest struct {
	Message string
}

// ChatEvent is one streamed event.
type ChatEvent struct {
	Type  string `json:"type"`
	Delta string `json:"delta,omitempty"`
}

// HealthResponse is a trivial health payload.
type HealthResponse struct {
	Status string `json:"status"`
}
