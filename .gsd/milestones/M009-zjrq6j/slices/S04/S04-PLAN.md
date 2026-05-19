# S04: Redis batch hit benchmark

**Goal:** Add benchmark sections for L1 hits, Redis L2 hits after restart, cached batch workloads, and repeated chunk reuse.
**Demo:** After this, benchmark evidence shows whether batch cache hits are Redis round-trip limited.

## Must-Haves

- Benchmark isolates L1 hot hit and Redis L2 after API restart.
- Benchmark includes cached batch workloads.
- Output includes p50/p95/p99 or equivalent latency summaries and Redis diagnostic deltas.
- Results can be compared with config snapshot from S01.

## Proof Level

- This slice proves: Python benchmark run with new sections and sanity checks

## Integration Closure

Provides measured go/no-go evidence for S05 MGET/pipeline implementation.

## Verification

- Adds Redis INFO deltas, hit/miss counts, latency percentiles, and workload labels to benchmark output.

## Tasks

- [x] **T01: Design batch hit benchmark sections** `est:small`
  Inspect current benchmark sections and batch API contract. Design minimal new benchmark helpers for batch requests, latency summaries, Redis INFO snapshots, and deltas without changing service behavior.
  - Files: `benchmark.py`, `api/handlers/batch.go`
  - Verify: Summary names endpoint payloads, Redis delta fields, and benchmark section layout.

- [x] **T02: Implement batch hit benchmark sections** `est:medium`
  Implement benchmark helper sections for L1 hot hit, Redis L2 after API restart, cached batch inputs, and repeated chunk reuse. Include Redis INFO deltas for each relevant section and avoid raw text output.
  - Files: `benchmark.py`
  - Verify: Python compile and parser checks for new sections pass.

- [x] **T03: Verify batch hit benchmark evidence** `est:medium`
  Run full benchmark verification with Docker stack, Go tests/lint, artifact parser for new sections and Redis deltas, and GitNexus detect_changes. Decide whether S05 should proceed or be skipped based on evidence.
  - Files: `benchmark.py`, `benchmark-results/`
  - Verify: `docker compose config`; `cd api && go test ./... -short`; pinned lint; `uv run --python 3.13 --with requests --with redis python benchmark.py`; parser checks; GitNexus detect_changes.

## Files Likely Touched

- benchmark.py
- api/handlers/batch.go
- benchmark-results/
