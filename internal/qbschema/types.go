// Package qbschema defines Go structs that mirror SigNoz's v5 Query Builder
// JSON format (the "compositeQuery" shape used by dashboards, alerts, and
// the query API). Filled in as the converter's target format is confirmed
// against SigNoz's docs and source.
package qbschema

// CompositeQuery is the root object SigNoz expects for a query.
type CompositeQuery struct {
	Queries []Query `json:"queries"`
}

// Query wraps a single builder query (later: also formula/promql/clickhouse_sql types).
type Query struct {
	Type string        `json:"type"` // "builder_query"
	Spec BuilderQuerySpec `json:"spec"`
}

// BuilderQuerySpec is the v5 builder_query spec.
type BuilderQuerySpec struct {
	Name         string        `json:"name"`
	Signal       string        `json:"signal"` // "metrics" | "logs" | "traces"
	Aggregations []Aggregation `json:"aggregations,omitempty"`
	Filter       *Filter       `json:"filter,omitempty"`
	GroupBy      []string      `json:"groupBy,omitempty"`
	Having       *Having       `json:"having,omitempty"`
	Limit        int           `json:"limit,omitempty"`
	OrderBy      string        `json:"orderBy,omitempty"`
}

// Aggregation is a single aggregation expression, e.g. "sum(value)".
type Aggregation struct {
	Expression string `json:"expression"`
}

// Filter is a SQL-like filter expression, e.g. "service = 'checkout'".
type Filter struct {
	Expression string `json:"expression"`
}

// Having is a post-aggregation filter expression, e.g. "count() > 100".
type Having struct {
	Expression string `json:"expression"`
}
