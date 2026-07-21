package convert

import (
	"fmt"

	"github.com/prometheus/prometheus/promql/parser"
)

// extractAggregation maps a PromQL aggregation op (sum, avg, count...)
// to a QB aggregation expression, e.g. sum -> sum(value).
//
// Note: doesn't unwrap nested calls yet (e.g. sum(rate(x[5m]))),  just reads the outer op

func extractAggregation(agg *parser.AggregateExpr) string {
	return fmt.Sprintf("%s(value)", agg.Op.String())
}