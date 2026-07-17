package httpbind

import (
	"net/http"
	"sync"
)

// OpenAPI document registration for generated embeds.
// Generation writes JSON/YAML bytes; handlers serve them without re-analyzing Go.

var (
	openAPIMu   sync.RWMutex
	openAPIJSON []byte
	openAPIYAML []byte
)

// RegisterOpenAPI stores generated OpenAPI document bytes for OpenAPIJSON/OpenAPIYAML.
// Called from generated init(); not a handwritten OpenAPI source of truth.
func RegisterOpenAPI(jsonDoc, yamlDoc []byte) {
	openAPIMu.Lock()
	defer openAPIMu.Unlock()
	openAPIJSON = append([]byte(nil), jsonDoc...)
	openAPIYAML = append([]byte(nil), yamlDoc...)
}

// OpenAPIDocumentJSON returns the registered OpenAPI JSON document (copy).
func OpenAPIDocumentJSON() []byte {
	openAPIMu.RLock()
	defer openAPIMu.RUnlock()
	if openAPIJSON == nil {
		return nil
	}
	return append([]byte(nil), openAPIJSON...)
}

// OpenAPIDocumentYAML returns the registered OpenAPI YAML document (copy).
func OpenAPIDocumentYAML() []byte {
	openAPIMu.RLock()
	defer openAPIMu.RUnlock()
	if openAPIYAML == nil {
		return nil
	}
	return append([]byte(nil), openAPIYAML...)
}

// OpenAPIJSON serves the embedded OpenAPI document as application/json.
func OpenAPIJSON(w http.ResponseWriter, r *http.Request) {
	_ = r
	doc := OpenAPIDocumentJSON()
	if len(doc) == 0 {
		WriteError(w, r, Internal(errNoOpenAPI))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(doc)
}

// OpenAPIYAML serves the embedded OpenAPI document as application/yaml.
func OpenAPIYAML(w http.ResponseWriter, r *http.Request) {
	_ = r
	doc := OpenAPIDocumentYAML()
	if len(doc) == 0 {
		WriteError(w, r, Internal(errNoOpenAPI))
		return
	}
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(doc)
}

type openAPIMissingError struct{}

func (openAPIMissingError) Error() string {
	return "httpbind: no OpenAPI document registered"
}

var errNoOpenAPI error = openAPIMissingError{}
