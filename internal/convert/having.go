package convert

import (
	"fmt"

	"github.com/prometheus/prometheus/promql/parser"
)

// extractHaving looks for a comparison wrapped around an aggregation,
// e.g. sum(errors) > 100, and turns it into a QB having expression
// like sum() > 100.

func extractHaving(bin *parser.BinaryExpr) (having string, agg *parser.AggregateExpr, ok bool) {
	aggExpr, isAgg := bin.LHS.(*parser.AggregateExpr)
	if !isAgg {
		return "", nil, false
	}

	num, isNum := bin.RHS.(*parser.NumberLiteral)
	if !isNum {
		return "", nil, false
	}

	op := binaryOpToOperator(bin.Op)
	return fmt.Sprintf("%s() %s %v", aggExpr.Op.String(), op, num.Val), aggExpr, true
}

func binaryOpToOperator(op parser.ItemType) string {
	switch op {
	case parser.GTR:
		return ">"
	case parser.LSS:
		return "<"
	case parser.GTE:
		return ">="
	case parser.LTE:
		return "<="
	case parser.EQLC:
		return "="
	case parser.NEQ:
		return "!="
	default:
		return "="
	}
}