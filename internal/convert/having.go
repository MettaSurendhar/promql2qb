package convert

import "github.com/prometheus/prometheus/promql/parser"

// extractHaving converts a post-aggregation comparison
// (e.g. sum(errors) > 100, expressed in PromQL as a BinaryExpr wrapping
// an AggregateExpr) into a QB having expression (e.g. "count() > 100").
//
// TODO: implement. Only needed when the top-level expr is a
// *parser.BinaryExpr whose LHS or RHS is an *parser.AggregateExpr.
func extractHaving(bin *parser.BinaryExpr) string {
	return ""
}
