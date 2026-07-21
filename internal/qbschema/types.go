package qbschema

type CompositeQuery struct {
	Queries []Query `json:"queries"`
}

type Query struct {
	Type string           `json:"type"` // "builder_query"
	Spec BuilderQuerySpec `json:"spec"`
}

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

type Aggregation struct {
	MetricName       string `json:"metricName"`
	TimeAggregation  string `json:"timeAggregation,omitempty"`  // rate | increase | latest
	SpaceAggregation string `json:"spaceAggregation"`           // sum | avg | count | min | max
}

type Filter struct {
	Expression string `json:"expression"`
}

type Having struct {
	Expression string `json:"expression"`
}