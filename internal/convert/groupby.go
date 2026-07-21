package convert

import "github.com/prometheus/prometheus/promql/parser"

// extractGroupBy pulls the "by (...)" labels off an aggregation.
// "without (...)" isn't handled yet
func extractGroupBy(agg *parser.AggregateExpr) []string {
	if agg.Without {
		return nil
	}
	return agg.Grouping
}