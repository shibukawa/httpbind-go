package httpbinder

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// SwaggerUI returns an http.Handler that serves a minimal Swagger UI page
// loading the OpenAPI document from specURL (e.g. "/openapi.json").
//
// Assets are loaded from a public CDN; this handler does not embed Swagger UI
// binaries. Mount freely, e.g.:
//
//	mux.Handle("GET /docs/{$}", httpbinder.SwaggerUI("/openapi.json"))
func SwaggerUI(specURL string) http.Handler {
	if strings.TrimSpace(specURL) == "" {
		specURL = "/openapi.json"
	}
	// Pre-quote for JS string literal; mark as template.JS to avoid double-escaping.
	specJS := template.JS(strconv.Quote(specURL))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = swaggerUITemplate.Execute(w, swaggerUIData{
			SpecURLJS: specJS,
			Title:     "httpbinder API docs",
			CSSURL:    swaggerUICSS,
			BundleJS:  swaggerUIBundleJS,
			PresetJS:  swaggerUIPresetJS,
		})
	})
}

// CDN pins (jsDelivr). Bump deliberately when upgrading Swagger UI.
const (
	swaggerUICSS      = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui.css"
	swaggerUIBundleJS = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"
	swaggerUIPresetJS = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui-standalone-preset.js"
)

type swaggerUIData struct {
	SpecURLJS template.JS
	Title     string
	CSSURL    string
	BundleJS  string
	PresetJS  string
}

var swaggerUITemplate = template.Must(template.New("swagger-ui").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>{{.Title}}</title>
  <link rel="stylesheet" href="{{.CSSURL}}"/>
  <style>
    body { margin: 0; background: #fafafa; }
    #swagger-ui { max-width: 100%; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="{{.BundleJS}}"></script>
  <script src="{{.PresetJS}}"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: {{.SpecURLJS}},
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
        layout: "StandaloneLayout",
        tryItOutEnabled: true
      });
    };
  </script>
</body>
</html>
`))
