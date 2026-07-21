package convert

import "github.com/prometheus/prometheus/promql/parser"

// extractFilter converts a PromQL vector selector's label matchers
// (e.g. {service="checkout", status!="200"}) into a QB filter expression
// (e.g. "service = 'checkout' AND status != '200'").
//
// TODO: implement. Walk sel.LabelMatchers, map each parser.MatchType
// (MatchEqual, MatchNotEqual, MatchRegexp, MatchNotRegexp) to the
// corresponding QB operator, and join with AND.
func extractFilter(sel *parser.VectorSelector) string {
	return ""
}
