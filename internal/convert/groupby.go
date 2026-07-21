package convert

import "github.com/prometheus/prometheus/promql/parser"

// extractGroupBy converts a PromQL "by (...)" / "without (...)" clause
// into a QB groupBy list.
//
// TODO: implement. agg.Grouping holds the label names; agg.Without
// distinguishes "by" from "without" (MVP: only support "by").
func extractGroupBy(agg *parser.AggregateExpr) []string {
	return nil
}
