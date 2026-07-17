package jsonbind

import (
	"strings"
	"testing"
)

type limitedDocument struct{ Value string }

func TestDecodeJSONLimitRejectsUnknownLengthReader(t *testing.T) {
	RegisterDecode[limitedDocument](func([]byte) (limitedDocument, error) { return limitedDocument{}, nil })
	_, err := DecodeJSONLimit[limitedDocument](strings.NewReader(`{"value":"too large"}`), 8)
	je, ok := AsError(err)
	if !ok || je.Code != "payload_too_large" {
		t.Fatalf("want payload_too_large, got %#v", err)
	}
}

type globallyLimitedDocument struct{ Value string }

func TestDecodeJSONUsesGlobalLimit(t *testing.T) {
	old := MaxJSONBodyBytes()
	SetMaxJSONBodyBytes(8)
	defer SetMaxJSONBodyBytes(old)
	RegisterDecode[globallyLimitedDocument](func([]byte) (globallyLimitedDocument, error) {
		return globallyLimitedDocument{}, nil
	})
	_, err := DecodeJSON[globallyLimitedDocument](strings.NewReader(`{"value":"too large"}`))
	je, ok := AsError(err)
	if !ok || je.Code != "payload_too_large" {
		t.Fatalf("want payload_too_large, got %#v", err)
	}
}
