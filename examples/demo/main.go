package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shibukawa/tinybind-go/configbind"
	"github.com/shibukawa/tinygodriver/httpmux"

	// Registers the host Netdever for TinyGo's net package.
	_ "github.com/shibukawa/tinygodriver/netdev"
)

func main() {
	cfg := configbind.Bind[ServerConfig]("server")
	if _, err := configbind.Load(configbind.LoadOptions{
		Vendor:   "shibukawa",
		Tool:     "tinybind-demo",
		FileName: "config.toml",
	}); err != nil {
		log.Fatalf("config: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	mux := httpmux.NewServeMux()
	RegisterDemoRoutes(mux)

	fmt.Printf("httpbind demo listening on http://localhost%s\n", addr)
	fmt.Printf("  docs:    http://localhost%s/docs/\n", addr)
	fmt.Printf("  openapi: http://localhost%s/openapi.json\n", addr)
	fmt.Printf("  index:   http://localhost%s/\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
