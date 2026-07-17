package sqlbind

import "testing"

func TestKeyAndConversions(t *testing.T) {
	r := Row{"id": int64(42), "name": []byte("alice"), "missing": nil}
	if k, p, e := Key(r, "id"); e != nil || !p || k != "42" {
		t.Fatalf("key %q %v %v", k, p, e)
	}
	if _, p, e := Key(r, "missing"); e != nil || p {
		t.Fatalf("null %v %v", p, e)
	}
	if v, e := Int(r, "id"); e != nil || v != 42 {
		t.Fatalf("int %v %v", v, e)
	}
	if v, e := String(r, "name"); e != nil || v != "alice" {
		t.Fatalf("string %q %v", v, e)
	}
}
