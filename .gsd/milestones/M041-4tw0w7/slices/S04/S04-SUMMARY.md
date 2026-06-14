---
id: S04
parent: M041-4tw0w7
milestone: M041-4tw0w7
provides:
  - Validated fd-controlled cache-hot performance envelope for downstream slices.
  - LRU cache implementation and cache observability for `/v1/embeddings`.
  - Explicit D045 semantics separating fd cache performance from backend miss inference latency.
requires:
  []
affects:
  []
key_files:
  - api/cache/lru.go
  - api/cache/lru_test.go
  - api/handlers/embeddings.go
  - api/fd_v2_cache_integration_test.go
  - api/observability/metrics.go
  - tools/verify_fd_v2_perf.sh
  - docs/fd-v2.md
  - benchmark-results/fd-v2-perf-validation-m041-s04.md
key_decisions:
  - D045: T-P-1..T-P-5 are cache-hot steady-state validation after prewarm; real cache-miss TEI CPU latency is diagnostic only for M041 S04.
patterns_established:
  - Use isolated Redis/cache namespace when comparing miss behavior, but use explicit prewarm and `X-Cache: HIT` assertions for cache-hot performance gates.
  - Keep cache-miss latency diagnostics visible even when not blocking the slice.
observability_surfaces:
  - `X-Cache: HIT|MISS` response header on `/v1/embeddings`.
  - `fd_cache_hits_total{result="hit|miss"}` and `fd_cache_evictions_total` Prometheus metrics.
  - Final performance artifact with cache-hot pass metrics and miss diagnostics.
drill_down_paths:
  - .gsd/milestones/M041-4tw0w7/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M041-4tw0w7/slices/S04/tasks/T02-SUMMARY.md
  - .gsd/milestones/M041-4tw0w7/slices/S04/tasks/T03-SUMMARY.md
  - .gsd/milestones/M041-4tw0w7/slices/S04/tasks/T04-SUMMARY.md
  - .gsd/milestones/M041-4tw0w7/slices/S04/tasks/T05-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-06-14T07:43:09.642Z
blocker_discovered: false
---

# S04: Performance baseline and LRU cache

**Delivered fd-controlled LRU caching, X-Cache/metrics visibility, and cache-hot performance validation while documenting TEI CPU miss latency as diagnostic.**

## What Happened

S04 established the performance baseline, implemented an in-memory LRU cache with TTL/size configuration, integrated cache status into `/v1/embeddings` responses via `X-Cache`, and added validation tooling for the accepted cache-hot T-P contract. The slice first surfaced an important reality: real cache-miss inference on the current TEI CPU backend misses the original latency targets and fd wrapper latency matches direct TEI latency. The user explicitly descoped backend remediation and chose D045: T-P-1..T-P-5 are cache-hot steady-state checks after explicit prewarm. The final verifier now prewarms measured payloads through real inference, requires `X-Cache: HIT` during measured latency cases, keeps cache-miss diagnostics in the artifact, and passes against the current fd service with real TEI/Redis.

## Verification

Fresh S04 final verification passed: `FD_BASE_URL=http://localhost:8000 tools/verify_fd_v2_perf.sh` exit 0. Artifact `benchmark-results/fd-v2-perf-validation-m041-s04.md` reports cache-hot batch=1 p95 2.236ms, batch=10 p95 3.468ms, batch=32 p95 7.595ms, 100 sequential cache-hot requests 0 errors/non-HIT, 4x8 concurrent cache-hot requests 0.010s, and repeated input HIT latency 1.870ms. Fresh M043 gates also passed after the final changes: `go test ./...` exit 0, golangci-lint v2.12.2 0 issues, govulncheck 0 reachable vulnerabilities.

## Requirements Advanced

- R012 — Clarified and validated as cache-hot steady-state performance after D045.
- R014 — Completed X-Cache header validation in addition to S03 headers.
- R016 — Implemented and validated LRU cache, hit/miss behavior, and metrics.

## Requirements Validated

- R012 — `benchmark-results/fd-v2-perf-validation-m041-s04.md` PASS under D045 cache-hot contract.
- R014 — Final verifier and integration tests prove `X-Cache: MISS|HIT`; S03 already validated remaining headers.
- R016 — LRU unit tests, cache integration test, metrics evidence, and final cache-hot runtime verifier.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

The original T-P wording was ambiguous and initially treated as real cache-miss inference. Real integration tests proved current TEI CPU miss latency cannot meet the targets. Per explicit user decision D045, T-P targets were rescoped to cache-hot steady-state; miss latency remains diagnostic only.

## Known Limitations

Real cache-miss TEI CPU inference is slow (diagnostic sample: batch=1 about 235ms, batch=10 about 2107ms, batch=32 about 6796ms) and is not remediated in M041 S04. Backend remediation such as ONNX/GPU/faster TEI runtime is explicitly out of scope for this slice.

## Follow-ups

If future work requires cache-miss latency targets, create a new milestone/slice for backend runtime remediation and legal-quality parity validation.

## Files Created/Modified

- `api/cache/lru.go` — In-memory LRU cache with TTL/size env config and EmbeddingCache adapter methods.
- `api/handlers/embeddings.go` — Adds `X-Cache` HIT/MISS status based on cache/model path.
- `api/fd_v2_cache_integration_test.go` — Integration-style cache MISS/HIT, latency, and metrics verification.
- `tools/verify_fd_v2_perf.sh` — Final cache-hot performance verifier with non-blocking miss diagnostics.
- `docs/fd-v2.md` — Clarifies R-P0-6 and T-P-1..T-P-5 cache-hot semantics.
