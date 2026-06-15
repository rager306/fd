---
estimated_steps: 1
estimated_files: 5
skills_used: []
---

# T02: Implemented health last_error/dependency/capacity fields and runtime/cache metrics gauges.

Extend health response with last_error, dependency reachability/latency, and in_flight_capacity. Add lightweight dependency checker abstractions and wire Redis/TEI probes in main. Add Prometheus gauges for in-flight requests/capacity and local cache entries, updated at scrape time or middleware time.

## Inputs

- `api/lifecycle/state.go`
- `api/cache/tiered.go`
- `api/embed/tei.go`

## Expected Output

- `api/handlers/health.go`
- `api/observability/metrics.go`
- `api/main.go`

## Verification

cd api && go test ./handlers ./observability && cd api && go test ./...

## Observability Impact

Adds durable diagnostic fields and gauges.
