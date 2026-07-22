package convert

import (
	"fmt"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql/parser"

	"github.com/MettaSurendhar/promql2qb/internal/qbschema"
)

func init() {
	model.NameValidationScheme = model.UTF8Validation
}

var promqlParser = parser.NewParser(parser.Options{})


func Convert(promql string) (*qbschema.CompositeQuery, error) {
	expr, err := promqlParser.ParseExpr(promql)
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

		aggregation, sel, err := extractAggregation(aggExpr)
		if err != nil {
			return nil, err
		}

		spec.Aggregations = []qbschema.Aggregation{aggregation}
		spec.GroupBy = extractGroupBy(aggExpr)
		spec.Having = &qbschema.Having{Expression: having}
		if f := extractFilter(sel); f != "" {
			spec.Filter = &qbschema.Filter{Expression: f}
		}

		return &qbschema.CompositeQuery{
			Queries: []qbschema.Query{{Type: "builder_query", Spec: spec}},
		}, nil
	}

	agg, ok := expr.(*parser.AggregateExpr)
	if !ok {
		return nil, fmt.Errorf("only aggregation queries are supported right now, e.g. sum(...) by (...)")
	}

	aggregation, sel, err := extractAggregation(agg)
	if err != nil {
		return nil, err
	}

	spec.Aggregations = []qbschema.Aggregation{aggregation}
	spec.GroupBy = extractGroupBy(agg)
	if f := extractFilter(sel); f != "" {
		spec.Filter = &qbschema.Filter{Expression: f}
	}

	return &qbschema.CompositeQuery{
		Queries: []qbschema.Query{{Type: "builder_query", Spec: spec}},
	}, nil
}