package httpbinder_test

import (
	"encoding/json"
	"testing"

	httpbinder "github.com/shibukawa/httpbind-go"
)

func TestRestJSONAny_ExcludesAndDecodesNested(t *testing.T) {
	body := map[string]json.RawMessage{
		"name":  json.RawMessage(`"Ada"`),
		"role":  json.RawMessage(`"admin"`),
		"meta":  json.RawMessage(`{"source":"import"}`),
		"count": json.RawMessage(`2`),
	}
	got, err := httpbinder.RestJSONAny(body, []string{"name", "email"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got["name"]; ok {
		t.Fatalf("name excluded: %#v", got)
	}
	if got["role"] != "admin" {
		t.Fatalf("role: %#v", got["role"])
	}
	meta, ok := got["meta"].(map[string]any)
	if !ok || meta["source"] != "import" {
		t.Fatalf("meta: %#v", got["meta"])
	}
	// JSON numbers decode as float64 into any
	if got["count"] != float64(2) {
		t.Fatalf("count: %#v (%T)", got["count"], got["count"])
	}
}

func TestRestJSONRaw_Copy(t *testing.T) {
	body := map[string]json.RawMessage{
		"name": json.RawMessage(`"x"`),
		"k":    json.RawMessage(`{"a":1}`),
	}
	got := httpbinder.RestJSONRaw(body, []string{"name"})
	if len(got) != 1 || string(got["k"]) != `{"a":1}` {
		t.Fatalf("%#v", got)
	}
}

func TestRestFormAny(t *testing.T) {
	got := httpbinder.RestFormAny(map[string]string{
		"name": "n",
		"x":    "1",
	}, []string{"name"})
	if got["x"] != "1" {
		t.Fatalf("%#v", got)
	}
	if _, ok := got["name"]; ok {
		t.Fatalf("name in rest: %#v", got)
	}
}
