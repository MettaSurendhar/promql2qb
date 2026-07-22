package convert

import (
	"encoding/json"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		promql  string
		wantErr bool
	}{
		{
			name:   "simple sum with label filter and group by",
			promql: `sum(http_requests_total{service="checkout"}) by (status)`,
		},
		{
			name:   "count with multiple label matchers",
			promql: `count(errors_total{service="checkout", status!="200"}) by (route)`,
		},
		{
			name:   "avg with regex matcher",
			promql: `avg(latency_seconds{service=~"checkout.*"}) by (route)`,
		},
		{
			name:   "having on a plain aggregation",
			promql: `sum(errors) > 100`,
		},
		{
			name:   "sum of rate over a range vector",
			promql: `sum(rate(http_requests_total{service="checkout"}[5m])) by (status)`,
		},
		{
			name:   "sum of increase over a range vector",
			promql: `sum(increase(errors_total{service="checkout"}[1m])) by (route)`,
		},
		{
			name:   "having on a rate aggregation",
			promql: `sum(rate(errors_total{service="checkout"}[5m])) > 10`,
		},
		{
			name:   "avg_over_time maps to avg time aggregation",
			promql: `sum(avg_over_time(gen{service="checkout"}[60s]))`,
		},
		{
			name:   "dotted label name stays quoted in the filter output",
			promql: `sum(avg_over_time(gen{"service.name"="telemetrygen"}[60s]))`,
		},
		{
			name:    "bare selector with no aggregation is not supported yet",
			promql:  `http_requests_total{service="checkout"}`,
			wantErr: true,
		},
		{
			name:    "binary expression between two metrics is not supported yet",
			promql:  `sum(errors) / sum(requests)`,
			wantErr: true,
		},
		{
			name:    "unsupported function inside the aggregation",
			promql:  `sum(abs(errors_total{service="checkout"})) by (route)`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := Convert(tt.promql)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Convert(%q): expected an error, got none", tt.promql)
				}
				t.Logf("got expected error: %v", err)
				return
			}

			if err != nil {
				t.Fatalf("Convert(%q): unexpected error: %v", tt.promql, err)
			}

			out, err := json.MarshalIndent(spec, "", "  ")
			if err != nil {
				t.Fatalf("Convert(%q): result did not marshal to JSON: %v", tt.promql, err)
			}
			t.Logf("input:  %s\noutput: %s", tt.promql, out)
		})
	}
}

// TestExtractFilterQuotesDottedLabelNames locks in the specific case that
// failed against a real SigNoz instance before the quoting fix.
func TestExtractFilterQuotesDottedLabelNames(t *testing.T) {
	spec, err := Convert(`sum(avg_over_time(gen{"service.name"="telemetrygen"}[60s])) by (status)`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := spec.Queries[0].Spec.Filter.Expression
	want := `"service.name" = 'telemetrygen'`
	if got != want {
		t.Fatalf("filter expression = %q, want %q", got, want)
	}
}