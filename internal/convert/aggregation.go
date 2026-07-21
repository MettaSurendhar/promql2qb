package convert

import (
	"fmt"

	"github.com/prometheus/prometheus/promql/parser"

	"github.com/MettaSurendhar/promql2qb/internal/qbschema"
)

// two aggregation layers:
// the outer op (sum/avg/count/min/max) becomes spaceAggregation,
// and an inner rate()/increase() call becomes timeAggregation.

func extractAggregation(agg *parser.AggregateExpr) (qbschema.Aggregation, *parser.VectorSelector, error) {
	metricName, sel, timeAgg, err := unwrapAggregatedExpr(agg.Expr)
	if err != nil {
		return qbschema.Aggregation{}, nil, err
	}

	return qbschema.Aggregation{
		MetricName:       metricName,
		TimeAggregation:  timeAgg,
		SpaceAggregation: agg.Op.String(),
	}, sel, nil
}

func unwrapAggregatedExpr(expr parser.Expr) (metricName string, sel *parser.VectorSelector, timeAgg string, err error) {
	switch e := expr.(type) {
	case *parser.VectorSelector:
		return e.Name, e, "latest", nil

	case *parser.Call:
		if len(e.Args) != 1 {
			return "", nil, "", fmt.Errorf("%s(...) with %d arguments isn't supported yet", e.Func.Name, len(e.Args))
		}

		matrixSel, ok := e.Args[0].(*parser.MatrixSelector)
		if !ok {
			return "", nil, "", fmt.Errorf("%s(...) needs a range vector like metric[5m], got something else", e.Func.Name)
		}
		vecSel, ok := matrixSel.VectorSelector.(*parser.VectorSelector)
		if !ok {
			return "", nil, "", fmt.Errorf("unsupported range selector inside %s(...)", e.Func.Name)
		}

		switch e.Func.Name {
		case "rate", "increase":
			return vecSel.Name, vecSel, e.Func.Name, nil
		default:
			return "", nil, "", fmt.Errorf("function %s(...) isn't supported yet", e.Func.Name)
		}

	default:
		return "", nil, "", fmt.Errorf("unsupported expression inside the aggregation")
	}
}