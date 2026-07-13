package mappingfixture_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawayoshiki/httpbind-go"
	"github.com/shibukawayoshiki/httpbind-go/generator"
	"github.com/shibukawayoshiki/httpbind-go/internal/mappingfixture"
)

func TestBind_JSONAndMetadata(t *testing.T) {
	body := `{"name":"Alice","email":"a@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/orgs/acme/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret")
	req.SetPathValue("org_id", "acme")

	got, err := httpbinder.Bind[mappingfixture.CreateUserRequest](req)
	if err != nil {
		t.Fatalf("Bind: %v", err)
	}
	if got.Name != "Alice" || got.Email != "a@example.com" {
		t.Fatalf("input fields: %+v", got)
	}
	if got.OrgID != "acme" {
		t.Fatalf("path org_id: %q", got.OrgID)
	}
	if got.Token != "Bearer secret" {
		t.Fatalf("header token: %q", got.Token)
	}
}

func TestBind_QueryInput(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users?name=Bob&email=b@example.com", nil)
	req.SetPathValue("org_id", "org1")
	req.Header.Set("Authorization", "t")

	got, err := httpbinder.Bind[mappingfixture.CreateUserRequest](req)
	if err != nil {
		t.Fatalf("Bind: %v", err)
	}
	if got.Name != "Bob" || got.Email != "b@example.com" {
		t.Fatalf("query input: %+v", got)
	}
	if got.OrgID != "org1" || got.Token != "t" {
		t.Fatalf("meta: %+v", got)
	}
}

func TestBind_SearchQueryAndPayload(t *testing.T) {
	body := `{"filter":"active"}`
	req := httptest.NewRequest(http.MethodPost, "/search?keyword=go&page=2", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	got, err := httpbinder.Bind[mappingfixture.SearchRequest](req)
	if err != nil {
		t.Fatalf("Bind: %v", err)
	}
	if got.Keyword != "go" || got.Page != 2 {
		t.Fatalf("query fields: %+v", got)
	}
	if got.Filter != "active" {
		t.Fatalf("payload filter: %q", got.Filter)
	}
}

func TestBind_FormPayload(t *testing.T) {
	form := "name=Carol&email=c@example.com"
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("org_id", "o")
	req.Header.Set("Authorization", "tok")

	got, err := httpbinder.Bind[mappingfixture.CreateUserRequest](req)
	if err != nil {
		t.Fatalf("Bind: %v", err)
	}
	if got.Name != "Carol" || got.Email != "c@example.com" {
		t.Fatalf("form bind: %+v", got)
	}
}

func TestWrite_JSONResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	err := httpbinder.Write[mappingfixture.CreateUserResponse](rec, req, mappingfixture.CreateUserResponse{
		ID:    "user_123",
		Name:  "Alice",
		Email: "a@example.com",
	})
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("content-type: %q", ct)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v\n%s", err, rec.Body.String())
	}
	if body["id"] != "user_123" || body["name"] != "Alice" || body["email"] != "a@example.com" {
		t.Fatalf("body: %#v", body)
	}
}

func TestWriteError_ValidationProblem(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	err := httpbinder.Validation(
		httpbinder.Field("email", "payload", "must be a valid email"),
	)
	httpbinder.WriteError(rec, req, err)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/problem+json") {
		t.Fatalf("content-type: %q", ct)
	}
	// Avoid map[string]any + interface type asserts here: TinyGo's encoding/json
	// can panic with reflect.AssignableTo when that path is linked with RawMessage bind.
	raw := rec.Body.String()
	if !strings.Contains(raw, `"status":400`) {
		t.Fatalf("status missing in %s", raw)
	}
	if !strings.Contains(raw, `"field":"email"`) || !strings.Contains(raw, `"location":"payload"`) {
		t.Fatalf("field error missing in %s", raw)
	}
}

func TestWriteError_HidesInternalCause(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	httpbinder.WriteError(rec, req, httpbinder.Internal(io.ErrUnexpectedEOF))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rec.Code)
	}
	raw := rec.Body.String()
	if strings.Contains(raw, "unexpected EOF") {
		t.Fatalf("internal cause leaked: %s", raw)
	}
}

func TestBind_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("org_id", "o")
	req.Header.Set("Authorization", "t")
	_, err := httpbinder.Bind[mappingfixture.CreateUserRequest](req)
	if err == nil {
		t.Fatal("expected error")
	}
	he, ok := httpbinder.AsHTTPError(err)
	if !ok || he.Status != http.StatusBadRequest {
		t.Fatalf("want 400 HTTPError, got %#v", err)
	}
}

func TestGenerator_EmitsTypeSpecificNoReflect(t *testing.T) {
	dir := t.TempDir()
	// copy types into temp package
	src, err := os.ReadFile("types.go")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "types.go"), src, 0o644); err != nil {
		t.Fatal(err)
	}
	// package name in types.go is mappingfixture — keep it
	out, err := generator.Generate(dir, dir, "httpbinder_gen.go")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	code := string(data)
	if !strings.Contains(code, "func bindCreateUserRequest") {
		t.Fatalf("missing bindCreateUserRequest in:\n%s", code)
	}
	if !strings.Contains(code, "func writeCreateUserResponse") {
		t.Fatalf("missing writeCreateUserResponse in:\n%s", code)
	}
	if !strings.Contains(code, "RegisterBind[CreateUserRequest]") {
		t.Fatalf("missing registration in:\n%s", code)
	}
	if strings.Contains(code, "\"reflect\"") || strings.Contains(code, "reflect.") {
		t.Fatalf("generated code must not use reflect:\n%s", code)
	}
	// field sources present as literals / calls
	for _, needle := range []string{
		`PathValue(r, "org_id")`,
		`HeaderValue(r, "Authorization")`,
		`QueryValue(r, "name")`,
		`QueryValue(r, "keyword")`,
	} {
		if !strings.Contains(code, needle) {
			t.Fatalf("missing %s in generated code", needle)
		}
	}
}

func TestGeneratedFile_NoReflectImport(t *testing.T) {
	data, err := os.ReadFile("httpbinder_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("\"reflect\"")) || bytes.Contains(data, []byte("reflect.")) {
		t.Fatal("committed generated file imports/uses reflect")
	}
}

func TestRoundTrip_HandlerStyle(t *testing.T) {
	// Real user path: Bind → service value → Write.
	// Call the handler directly (not via Go 1.22 method-path ServeMux patterns),
	// so TinyGo's net/http (without full pattern routing) can exercise the same I/O.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input, err := httpbinder.Bind[mappingfixture.CreateUserRequest](r)
		if err != nil {
			httpbinder.WriteError(w, r, err)
			return
		}
		out := mappingfixture.CreateUserResponse{
			ID:    "user_123",
			Name:  input.Name,
			Email: input.Email,
		}
		if err := httpbinder.Write[mappingfixture.CreateUserResponse](w, r, out); err != nil {
			httpbinder.WriteError(w, r, err)
		}
	})

	body := `{"name":"Alice","email":"a@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/orgs/acme/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer x")
	req.SetPathValue("org_id", "acme")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
	var m map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil {
		t.Fatal(err)
	}
	if m["id"] != "user_123" || m["name"] != "Alice" || m["email"] != "a@example.com" {
		t.Fatalf("response: %#v", m)
	}
}
