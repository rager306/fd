# S02: Async parallel chunked TEI calls in handler

**Goal:** Async parallel chunked TEI calls in handler with M043 lint/security gates built into the design: bounded concurrency, context propagation, small helper functions, observability, perf proof, and regression coverage.
**Demo:** After this, FD_ASYNC_CHUNKS=true enables parallel chunking: cold path for batch=128 падает с 25s до ≤10s, batch=32 с 6s до ≤4s. Cache hit path не regressed. New X-Concurrent-Chunks header в response для observability.

## Must-Haves

- FD_ASYNC_CHUNKS=true enables bounded parallel chunking with max concurrency 4.
- Cold path batch=128 improves to <=10s without cache-hit regression.
- X-Concurrent-Chunks header or equivalent metric exposes concurrency count in async mode.
- Production functions introduced/modified by this slice stay at gocyclo <=15 or have explicit justified exceptions.
- `go test ./...`, 18-linter golangci-lint, and govulncheck all exit 0.

## Proof Level

- This slice proves: Perf proof + integration/regression proof + M043 static-analysis gates.

## Integration Closure

Async mode remains opt-in. Default sync TEI path is unchanged and regression-tested. M041 acceptance behavior remains valid in both FD_ASYNC_CHUNKS=false and true modes.

## Verification

- Adds async chunk count and error/cancellation visibility; must avoid high-cardinality metric labels and noisy global state.

## Tasks

- [ ] **T01: Implement lint-aware async chunked orchestrator** `est:3h`
  Implement bounded parallel TEI chunk orchestration for batches larger than TEI's per-request limit. Keep production functions below gocyclo 15 by extracting small helpers for chunk planning, worker execution, ordered result assembly, and error aggregation. Propagate request context into every goroutine/call so contextcheck stays clean; ensure goroutines stop on cancellation and do not leak. Avoid exported APIs unless required; any exported helper/type needs meaningful godoc because revive:exported is now enforced.
  - Files: `api/handlers/embeddings.go`, `api/handlers/embeddings_integration_test.go`
  - Verify: cd api && go test ./handlers ./middleware && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./handlers ./middleware

- [ ] **T02: Wire FD_ASYNC_CHUNKS env into handler and main config** `est:1h`
  Add FD_ASYNC_CHUNKS configuration with default false and wire it into EmbeddingsHandler without changing default TEI behavior. Keep config parsing small and testable; use contextual errors for invalid env values if any. Ensure the handler path is explicit enough for future agents to see sync vs async behavior, and keep public config comments/godoc aligned with revive:exported requirements if new exported symbols are introduced.
  - Files: `api/main.go`, `api/handlers/embeddings.go`, `api/main_test.go`, `api/handlers/embeddings_integration_test.go`
  - Verify: cd api && go test ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...

- [ ] **T03: Add async observability header and metrics** `est:1h`
  Add X-Concurrent-Chunks response header in async mode and metrics/log surfaces that make the chunk count, miss count, and cancellation/error paths observable. Do not introduce noisy globals or unbounded labels. Ensure header/middleware interactions remain compatible with M041 S03 headers work and that new observability code passes gocritic/contextcheck.
  - Files: `api/handlers/embeddings.go`, `api/handlers/embeddings_integration_test.go`
  - Verify: cd api && go test ./handlers && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./handlers

- [ ] **T04: Benchmark async versus sync cold and warm paths** `est:2h`
  Create or update perf benchmark tooling to compare FD_ASYNC_CHUNKS=false versus true for cold batch=32, cold batch=128, and cache-hit path. Record before/after numbers in benchmark-results/fd-v2-async-perf-m042.md. Include the expanded M043 gate in the verification transcript so perf work cannot regress lint/test/security gates.
  - Files: `tools/verify_fd_async_perf.sh`, `benchmark-results/fd-v2-async-perf-m042.md`, `benchmark.py`
  - Verify: cd api && go test ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./... && go run golang.org/x/vuln/cmd/govulncheck@latest ./... && ../tools/verify_fd_async_perf.sh

- [ ] **T05: Run regression suite for M041 acceptance in async mode** `est:1h`
  Run the M041 acceptance/regression suite with FD_ASYNC_CHUNKS=false and true. Confirm error envelopes, validation behavior, cache-hit path, headers, and lifecycle assumptions are unchanged. Final completion evidence must include go test ./..., golangci-lint 18 linters, and govulncheck 0 reachable vulnerabilities.
  - Files: `benchmark-results/fd-v2-async-regression-m042.md`, `api/handlers/embeddings_integration_test.go`
  - Verify: cd api && go test ./... && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./... && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## Files Likely Touched

- api/handlers/embeddings.go
- api/handlers/embeddings_integration_test.go
- api/main.go
- api/main_test.go
- tools/verify_fd_async_perf.sh
- benchmark-results/fd-v2-async-perf-m042.md
- benchmark.py
- benchmark-results/fd-v2-async-regression-m042.md
