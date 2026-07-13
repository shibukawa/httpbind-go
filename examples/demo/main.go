package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}

	mux := http.NewServeMux()
	RegisterDemoRoutes(mux)

	fmt.Printf("httpbinder demo listening on http://localhost%s\n", addr)
	fmt.Printf("  docs:    http://localhost%s/docs/\n", addr)
	fmt.Printf("  openapi: http://localhost%s/openapi.json\n", addr)
	fmt.Printf("  index:   http://localhost%s/\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
