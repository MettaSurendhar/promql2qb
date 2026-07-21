// Command promql2qb reads a PromQL query and prints the equivalent
// SigNoz Query Builder JSON.
//
// Usage:
//
//	promql2qb 'sum(http_requests_total{service="checkout"}) by (status)'
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MettaSurendhar/promql2qb/internal/convert"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: promql2qb '<promql query>'")
		os.Exit(1)
	}

	spec, err := convert.Convert(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}
