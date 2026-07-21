// Package convert turns a PromQL query string into a SigNoz Query Builder
// (qbschema.CompositeQuery) value.
//
// Pipeline: PromQL string -> promql/parser AST -> Convert() walks the AST
// and delegates to filter.go, aggregation.go, groupby.go, having.go to
// build up the target spec.
package convert

import (
	"fmt"

	"github.com/prometheus/prometheus/promql/parser"

	"github.com/MettaSurendhar/promql2qb/internal/qbschema"
)

// Convert parses a PromQL query string and returns the equivalent SigNoz
// Query Builder composite query.
func Convert(promql string) (*qbschema.CompositeQuery, error) {
	expr, err := parser.ParseExpr(promql)
	if err != nil {
		return nil, fmt.Errorf("parsing promql: %w", err)
	}

	spec := qbschema.BuilderQuerySpec{
		Name:   "A",
		Signal: "metrics",
	}

	// TODO: walk expr (a parser.Expr) with a type switch and fill in spec
	// via extractFilter, extractAggregation, extractGroupBy, extractHaving.
	_ = expr

	return &qbschema.CompositeQuery{
		Queries: []qbschema.Query{
			{Type: "builder_query", Spec: spec},
		},
	}, fmt.Errorf("not implemented yet")
}
