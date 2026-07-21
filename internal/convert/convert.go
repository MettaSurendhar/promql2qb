package convert

import (
	"fmt"

	"github.com/prometheus/prometheus/promql/parser"

	"github.com/MettaSurendhar/promql2qb/internal/qbschema"
)

// Convert parses a PromQL query and builds the matching SigNoz
// Query Builder spec.
func Convert(promql string) (*qbschema.CompositeQuery, error) {
	expr, err := parser.ParseExpr(promql)
	if err != nil {
		return nil, fmt.Errorf("parsing promql: %w", err)
	}

	spec := qbschema.BuilderQuerySpec{
		Name:   "A",
		Signal: "metrics",
	}

	if bin, isBin := expr.(*parser.BinaryExpr); isBin {
		having, aggExpr, ok := extractHaving(bin)
		if !ok {
			return nil, fmt.Errorf("unsupported comparison shape for having")
		}

		spec.Aggregations = []qbschema.Aggregation{{Expression: extractAggregation(aggExpr)}}
		spec.GroupBy = extractGroupBy(aggExpr)
		spec.Having = &qbschema.Having{Expression: having}

		if sel, ok := aggExpr.Expr.(*parser.VectorSelector); ok {
			if f := extractFilter(sel); f != "" {
				spec.Filter = &qbschema.Filter{Expression: f}
			}
		}

		return &qbschema.CompositeQuery{
			Queries: []qbschema.Query{{Type: "builder_query", Spec: spec}},
		}, nil
	}

	agg, ok := expr.(*parser.AggregateExpr)
	if !ok {
		return nil, fmt.Errorf("only aggregation queries are supported right now, e.g. sum(...) by (...)")
	}

	spec.Aggregations = []qbschema.Aggregation{
		{Expression: extractAggregation(agg)},
	}
	spec.GroupBy = extractGroupBy(agg)

	sel, ok := agg.Expr.(*parser.VectorSelector)
	if !ok {
		return nil, fmt.Errorf("only a plain metric selector inside the aggregation is supported right now")
	}

	if f := extractFilter(sel); f != "" {
		spec.Filter = &qbschema.Filter{Expression: f}
	}

	return &qbschema.CompositeQuery{
		Queries: []qbschema.Query{
			{Type: "builder_query", Spec: spec},
		},
	}, nil
}