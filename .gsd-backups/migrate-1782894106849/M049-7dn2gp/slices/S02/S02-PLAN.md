# S02: Health and metrics context

**Goal:** Implement issue #8 AN-B/AN-C: richer health context and cheap Prometheus capacity/cache gauges for agent diagnostics.
**Demo:** After this slice, /health and /metrics expose last error, dependency, capacity, and cache occupancy signals agents can inspect.

## Must-Haves

- /health includes `last_error` when lifecycle state is not ok and an error exists.
- /health includes lightweight `dependencies.tei` and `dependencies.redis` reachability/latency metadata.
- /health includes `in_flight_capacity` alongside `in_flight_requests`.
- /metrics exposes in-flight and L1 cache occupancy gauges.
- Tests prove fields and metrics output without requiring live external services.

## Proof Level

- This slice proves: Red/green health and metrics tests plus full Go tests; live container smoke in S03.

## Integration Closure

Health dependency probes must be lightweight connectivity probes, not embedding inference; readiness semantics remain unchanged.

## Verification

- Agents can diagnose warmup/runtime failures, dependency reachability, capacity pressure, and cache occupancy from HTTP surfaces.

## Tasks

- [x] **T01: Pinned health/metrics diagnostic gaps with red tests.** `est:medium`
  Add failing tests for last_error serialization, dependency blocks, in_flight_capacity, and Prometheus gauges for in-flight/cache occupancy.
  - Files: `api/handlers/health_test.go`, `api/observability/metrics_test.go`, `api/main_test.go`
  - Verify: cd api && go test ./handlers ./observability (expected red before implementation).

- [x] **T02: Implemented health last_error/dependency/capacity fields and runtime/cache metrics gauges.** `est:large`
  Extend health response with last_error, dependency reachability/latency, and in_flight_capacity. Add lightweight dependency checker abstractions and wire Redis/TEI probes in main. Add Prometheus gauges for in-flight requests/capacity and local cache entries, updated at scrape time or middleware time.
  - Files: `api/handlers/health.go`, `api/handlers/health_test.go`, `api/observability/metrics.go`, `api/observability/metrics_test.go`, `api/main.go`
  - Verify: cd api && go test ./handlers ./observability && cd api && go test ./...

- [x] **T03: Recorded S02 evidence, advanced R041, and saved UAT for AN-B/AN-C.** `est:small`
  Write S02 evidence artifact, run static proof for AN-B/AN-C, update R041 with source/test proof while leaving live container proof for S03, save UAT, and complete S02.
  - Files: `benchmark-results/m049-s02-health-metrics-context.md`, `.gsd/REQUIREMENTS.md`
  - Verify: Artifact/static UAT proves health/metrics source and tests pass.

## Files Likely Touched

- api/handlers/health_test.go
- api/observability/metrics_test.go
- api/main_test.go
- api/handlers/health.go
- api/observability/metrics.go
- api/main.go
- benchmark-results/m049-s02-health-metrics-context.md
- .gsd/REQUIREMENTS.md
