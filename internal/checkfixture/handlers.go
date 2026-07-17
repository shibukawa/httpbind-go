package checkfixture

import (
	"net/http"

	httpbind "github.com/shibukawa/tinybind-go"
)

// RegisterRoutes exposes OpenAPICheck for OpenAPI generation tests.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /check", handleOpenAPICheck)
}

func handleOpenAPICheck(w http.ResponseWriter, r *http.Request) {
	in, err := httpbind.Bind[OpenAPICheckRequest](r)
	if err != nil {
		httpbind.WriteError(w, r, err)
		return
	}
	_ = in
	if err := httpbind.Write[OpenAPICheckResponse](w, r, OpenAPICheckResponse{OK: true}); err != nil {
		httpbind.WriteError(w, r, err)
	}
}
