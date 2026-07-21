package main

import (
	"bytes"
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

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	if err := enc.Encode(spec); err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling output: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(buf.String())
}