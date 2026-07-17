package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shibukawa/httpbind-go/generator"
	"github.com/shibukawa/httpbind-go/parser"
)

func main() {
	dir := flag.String("dir", ".", "package directory to analyze")
	out := flag.String("out", "", "output directory (default: same as -dir)")
	name := flag.String("name", "httpbinder_gen.go", "binder/writer output file name")
	openapi := flag.Bool("openapi", true, "also generate OpenAPI embed (httpbinder_openapi_gen.go)")
	openapiName := flag.String("openapi-name", "httpbinder_openapi_gen.go", "OpenAPI output file name")
	check := flag.Bool("check", false, "report analysis diagnostics and exit 1 if any undiscoverable route candidates exist")
	generateAll := flag.Bool("generate-all", false, "generate every mapping path for every struct (legacy mode)")
	flag.Parse()

	if *check {
		diags, err := parser.CheckPackage(*dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "httpbinder-gen check: %v\n", err)
			os.Exit(1)
		}
		for _, d := range diags {
			fmt.Fprintln(os.Stderr, d.String())
		}
		if len(diags) > 0 {
			fmt.Fprintf(os.Stderr, "httpbinder-gen check: %d diagnostic(s)\n", len(diags))
			os.Exit(1)
		}
		fmt.Println("ok")
		return
	}

	path, err := generator.New(generator.Options{GenerateAll: *generateAll}).Generate(*dir, *out, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpbinder-gen: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(path)

	if *openapi {
		op, err := generator.GenerateOpenAPI(*dir, *out, *openapiName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "httpbinder-gen openapi: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(op)
	}
}
