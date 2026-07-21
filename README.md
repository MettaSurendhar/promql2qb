# promql2qb

Convert PromQL queries into [SigNoz](https://github.com/SigNoz/signoz) Query Builder (v5) JSON — so people migrating from Prometheus/Grafana can bring their existing queries into SigNoz instead of rebuilding them by hand.

Built for the [Agents of SigNoz](https://www.wemakedevs.org/hackathons/signoz) hackathon (SigNoz x WeMakeDevs), Track 2 — Signals & Dashboards.

## Why

SigNoz's Query Builder is expression-based and can express most of what PromQL expresses, but there's no automated way to go from a PromQL string to the equivalent Query Builder JSON. This project fills that gap with a small, testable Go library + CLI.

## Status

🚧 Early scaffold — architecture and types are in place, conversion logic is not implemented yet. See [Roadmap](#roadmap).

## How it works

```
PromQL string
      │
      ▼
Parser (prometheus/promql/parser)
      │
      ▼
Internal AST (selectors, matchers, functions)
      │
      ▼
Mapper (AST node -> QB expressions)
      │
      ▼
SigNoz Query Builder JSON (compositeQuery spec)
```

The parser is not hand-written — it reuses Prometheus's own `promql/parser` package, so PromQL support is spec-correct from day one. The mapper walks the resulting AST and builds up a `qbschema.CompositeQuery` value, which is then marshalled to JSON.

## Project layout

```
promql2qb/
├── cmd/promql2qb/main.go        CLI entrypoint: reads a query, prints JSON
├── internal/
│   ├── convert/
│   │   ├── convert.go           top-level Convert(promqlString) (*qbschema.CompositeQuery, error)
│   │   ├── filter.go            label matchers -> filter.expression
│   │   ├── aggregation.go       agg functions -> aggregations[]
│   │   ├── groupby.go           by/without -> groupBy
│   │   ├── having.go            post-aggregation comparisons -> having.expression
│   │   └── convert_test.go      table-driven tests
│   └── qbschema/
│       └── types.go             Go structs mirroring SigNoz's v5 builder_query JSON
├── examples/                    PromQL -> JSON example pairs (used as docs + test fixtures)
├── go.mod
└── README.md
```

## MVP scope

**In scope:**
- Simple vector selectors with label matchers, e.g. `http_requests_total{service="checkout", status!="200"}`
- One aggregation at a time: `sum`, `avg`, `count`, `rate`
- `by (label)` group-by
- Post-aggregation comparisons (`sum(...) > 100`) mapped to `having`

**Out of scope for now** (tracked as future work):
- Nested/binary PromQL expressions (`a / b`, subqueries)
- `histogram_quantile`, `topk` / `bottomk`
- `without (...)` grouping
- Metrics-specific temporal/spatial aggregation split

## Usage (once implemented)

```bash
go build -o promql2qb ./cmd/promql2qb
./promql2qb 'sum(http_requests_total{service="checkout"}) by (status)'
```

```json
{
  "queries": [
    {
      "type": "builder_query",
      "spec": {
        "name": "A",
        "signal": "metrics",
        "aggregations": [{ "expression": "sum(value)" }],
        "filter": { "expression": "service = 'checkout'" },
        "groupBy": ["status"]
      }
    }
  ]
}
```

## Development

```bash
go mod tidy
go test ./...
```

## Roadmap

- [ ] Implement `extractFilter` (label matchers -> filter expression)
- [ ] Implement `extractAggregation` (agg functions -> aggregation expression)
- [ ] Implement `extractGroupBy` (`by (...)`)
- [ ] Implement `extractHaving` (post-aggregation comparisons)
- [ ] Wire it all up in `Convert`
- [ ] Add example pairs + table-driven test cases per feature
- [ ] Validate output against a self-hosted SigNoz instance
- [ ] README compatibility matrix (✅ / ❌ per PromQL feature)

## License

MIT
