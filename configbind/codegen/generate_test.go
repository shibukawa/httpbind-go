package codegen

import (
	"strings"
	"testing"

	"github.com/shibukawa/tinybind-go/minitoml"
)

func TestGenerateScaffolds(t *testing.T) {
	tomlText, envText, err := GenerateScaffolds([]Spec{
		{
			TypeName: "WebServerConfig",
			Prefix:   "webserver",
			Fields: []Field{
				{GoName: "Port", Key: "port", Kind: FieldInt, Default: "8080", Opt: "port,p", Help: "HTTP listen port"},
				{GoName: "Host", Key: "host", Kind: FieldString, Default: "localhost", Help: "listen host"},
				{GoName: "Origins", Key: "origins", Kind: FieldStringSlice},
				{GoName: "Secret", Key: "secret", Kind: FieldString, Env: "-"},
				{GoName: "TLS", Key: "tls", Kind: FieldStruct, Nested: []Field{
					{GoName: "Enabled", Key: "enabled", Kind: FieldBool, Default: "true", Help: "enable TLS"},
				}},
			},
		},
		{
			TypeName: "CacheConfig",
			Prefix:   "middleware.cache",
			Fields: []Field{
				{GoName: "ServiceName", Key: "service_name", Kind: FieldString, Env: "OTEL_SERVICE_NAME"},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	wantTOML := `[webserver]
# HTTP listen port
port = 8080
# listen host
host = "localhost"
origins = []
secret = ""
# enable TLS
tls.enabled = true

[middleware.cache]
service_name = ""
`
	if tomlText != wantTOML {
		t.Fatalf("TOML scaffold:\n--- got ---\n%s--- want ---\n%s", tomlText, wantTOML)
	}
	if _, err := minitoml.ParseString(tomlText); err != nil {
		t.Fatalf("generated TOML does not parse: %v\n%s", err, tomlText)
	}
	wantEnv := `# HTTP listen port
PORT=8080
# listen host
WEBSERVER_HOST="localhost"
WEBSERVER_ORIGINS=""
# enable TLS
WEBSERVER_TLS_ENABLED=true
OTEL_SERVICE_NAME=""
`
	if envText != wantEnv {
		t.Fatalf("env scaffold:\n--- got ---\n%s--- want ---\n%s", envText, wantEnv)
	}
}

func TestGenerateEmitsScaffoldConstants(t *testing.T) {
	src, err := Generate("fixture", []Spec{{
		TypeName: "ServerConfig",
		Prefix:   "server",
		Fields:   []Field{{GoName: "Port", Key: "port", Kind: FieldInt, Default: "8080"}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`const ConfigbindScaffoldTOML = "[server]\nport = 8080\n"`,
		`const ConfigbindScaffoldEnv = "SERVER_PORT=8080\n"`,
	} {
		if !strings.Contains(string(src), want) {
			t.Fatalf("generated scaffold constant %q missing:\n%s", want, src)
		}
	}
}

func TestGenerateEmitsEnvironmentOverride(t *testing.T) {
	src, err := Generate("fixture", []Spec{{
		TypeName: "ObservabilityConfig",
		Prefix:   "observability",
		Fields: []Field{
			{GoName: "ServiceName", Key: "service_name", Kind: FieldString, Env: "OTEL_SERVICE_NAME"},
		},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(src), `Env: "OTEL_SERVICE_NAME"`) {
		t.Fatalf("generated environment override missing:\n%s", src)
	}
}

func TestGenerateRejectsDuplicateEnvironmentOverride(t *testing.T) {
	_, err := Generate("fixture", []Spec{{
		TypeName: "ObservabilityConfig",
		Prefix:   "observability",
		Fields: []Field{
			{GoName: "ServiceName", Key: "service_name", Kind: FieldString, Env: "OTEL_SERVICE_NAME"},
			{GoName: "PeerName", Key: "peer_name", Kind: FieldString, Env: "OTEL_SERVICE_NAME"},
		},
	}})
	if err == nil || !strings.Contains(err.Error(), "duplicate environment variable") {
		t.Fatalf("error=%v", err)
	}
}
