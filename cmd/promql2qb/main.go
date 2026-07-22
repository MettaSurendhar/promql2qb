package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MettaSurendhar/promql2qb/internal/convert"
	"github.com/MettaSurendhar/promql2qb/internal/qbschema"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "promql2qb: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	query, err := readQuery()
	if err != nil {
		return err
	}

	spec, err := convert.Convert(query)
	if err != nil {
		return fmt.Errorf("converting query: %w", err)
	}

	out, err := formatJSON(spec)
	if err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}

	fmt.Printf("query:  %s\n\n%s\n", query, out)
	return nil
}

func readQuery() (string, error) {
	switch len(os.Args) {
	case 1:
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("reading query from stdin: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			return "", fmt.Errorf("no query given: pass it as an argument or pipe it via stdin")
		}
		return line, nil
	case 2:
		return os.Args[1], nil
	default:
		return "", fmt.Errorf("usage: promql2qb '<promql query>'  (or pipe it via stdin)")
	}
}

func formatJSON(spec *qbschema.CompositeQuery) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(spec); err != nil {
		return "", err
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}