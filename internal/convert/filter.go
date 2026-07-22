package convert

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

var validIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// extractFilter builds a QB filter expression from a vector selector's
// label matchers, e.g. {service="checkout", status!="200"} becomes
// service = 'checkout' AND status != '200'.

func extractFilter(sel *parser.VectorSelector) string {
	var parts []string

	for _, m := range sel.LabelMatchers {
		if m.Name == "__name__" {
			continue
		}

		name := m.Name
		if !validIdentifier.MatchString(name) {
			name = fmt.Sprintf("%q", name)
		}

		op := matchTypeToOperator(m.Type)
		parts = append(parts, fmt.Sprintf("%s %s '%s'", name, op, m.Value))
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