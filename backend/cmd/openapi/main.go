// Command openapi generates and exports the Wearlog OpenAPI 3.1 specification.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/StephenQiu30/wearlog-server/backend/internal/adapter/httpapi"
)

func main() {
	outPath := flag.String("out", "", "Output path for the generated openapi.json (default writes to stdout)")
	format := flag.String("format", "json", "Output format: json or yaml")
	flag.Parse()

	ctx := context.Background()
	yamlDoc, jsonDoc, err := httpapi.GeneratedOpenAPI(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate openapi schema: %v\n", err)
		os.Exit(1)
	}

	var data []byte
	switch *format {
	case "yaml", "yml":
		data = yamlDoc
	case "json":
		data = jsonDoc
	default:
		fmt.Fprintf(os.Stderr, "unsupported format %q, use json or yaml\n", *format)
		os.Exit(1)
	}

	if *outPath == "" {
		if _, err := os.Stdout.Write(data); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to stdout: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := os.MkdirAll(filepath.Dir(*outPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create directory for %s: %v\n", *outPath, err)
		os.Exit(1)
	}

	if err := os.WriteFile(*outPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write file %s: %v\n", *outPath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully exported OpenAPI spec (%s) to %s (%d bytes)\n", *format, *outPath, len(data))
}
