package convert

import "testing"

// Table-driven tests, one row per PromQL feature. Add cases as each
// piece of the converter is implemented; keep the JSON fixtures in
// examples/ once there are enough to warrant separate files.
func TestConvert(t *testing.T) {
	tests := []struct {
		name   string
		promql string
	}{
		{
			name:   "simple sum with label filter and group by",
			promql: `sum(http_requests_total{service="checkout"}) by (status)`,
		},
		// TODO: add cases for count, avg, rate, having, without, etc.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("not implemented yet")
			_, err := Convert(tt.promql)
			if err != nil {
				t.Fatalf("Convert(%q) error: %v", tt.promql, err)
			}
		})
	}
}
