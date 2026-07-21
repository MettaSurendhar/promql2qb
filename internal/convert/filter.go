package convert

import (
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

// extractFilter builds a QB filter expression from a vector selector's
// label matchers, e.g. {service="checkout", status!="200"} becomes
// service = 'checkout' AND status != '200'.
func extractFilter(sel *parser.VectorSelector) string {
	var parts []string

	for _, m := range sel.LabelMatchers {
		if m.Name == "__name__" {
			// this is the metric name itself, not a label filter
			continue
		}

		op := matchTypeToOperator(m.Type)
		parts = append(parts, fmt.Sprintf("%s %s '%s'", m.Name, op, m.Value))
	}

	return strings.Join(parts, " AND ")
}

func matchTypeToOperator(t labels.MatchType) string {
	switch t {
	case labels.MatchEqual:
		return "="
	case labels.MatchNotEqual:
		return "!="
	case labels.MatchRegexp:
		return "=~"
	case labels.MatchNotRegexp:
		return "!~"
	default:
		return "="
	}
}