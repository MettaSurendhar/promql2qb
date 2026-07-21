package convert

import "github.com/prometheus/prometheus/promql/parser"

// extractAggregation converts a PromQL aggregation expression
// (e.g. sum(rate(metric[5m]))) into a QB aggregation expression
// (e.g. "sum(value)").
//
// TODO: implement. Handle at least sum, avg, count, min, max, rate.
// AggregateExpr.Op is the aggregation function; AggregateExpr.Expr is
// the inner expression (which may itself be a Call, e.g. rate(...)).
func extractAggregation(agg *parser.AggregateExpr) string {
	return ""
}
