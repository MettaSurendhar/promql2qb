package convert

import (
	"fmt"

	"github.com/prometheus/prometheus/promql/parser"

	"github.com/MettaSurendhar/promql2qb/internal/qbschema"
)


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

var overTimeToTimeAgg = map[string]string{
	"avg_over_time":   "avg",
	"sum_over_time":   "sum",
	"min_over_time":   "min",
	"max_over_time":   "max",
	"count_over_time": "count",
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

		if e.Func.Name == "rate" || e.Func.Name == "increase" {
			return vecSel.Name, vecSel, e.Func.Name, nil
		}
		if timeAgg, ok := overTimeToTimeAgg[e.Func.Name]; ok {
			return vecSel.Name, vecSel, timeAgg, nil
		}
		return "", nil, "", fmt.Errorf("function %s(...) isn't supported yet", e.Func.Name)

	default:
		return "", nil, "", fmt.Errorf("unsupported expression inside the aggregation")
	}
}