# S04: Performance baseline and LRU cache — UAT

**Milestone:** M041-4tw0w7
**Written:** 2026-06-14T07:43:09.642Z

# S04 UAT — Performance baseline and LRU cache

## Result
PASS under D045 cache-hot steady-state contract.

## Checks
- UAT-01 LRU behavior: PASS. Unit and integration tests cover copy safety, dimension-specific keys, TTL expiry, eviction, concurrency, and first MISS / repeat HIT behavior.
- UAT-02 X-Cache observability: PASS. `/v1/embeddings` emits `X-Cache: MISS` for prewarm/cache-fill and `X-Cache: HIT` for repeated payloads.
- UAT-03 Metrics: PASS. Cache hit/miss counters are exposed and verified by integration tests.
- UAT-04 Final performance validation: PASS. `tools/verify_fd_v2_perf.sh` validates cache-hot T-P-1..T-P-5 against current fd + real TEI/Redis and writes `benchmark-results/fd-v2-perf-validation-m041-s04.md`.
- UAT-05 Real cache-miss transparency: PASS. The final artifact includes non-blocking cache-miss diagnostics showing TEI CPU miss latency remains slow and out of fd S04 scope.

## Evidence
- `benchmark-results/fd-v2-perf-validation-m041-s04.md`
- `benchmark-results/m041-s04-t05-cache-hot-runtime.txt`
- `benchmark-results/m041-s04-t05-go-test.txt`
- `benchmark-results/m041-s04-t05-lint.txt`
- `benchmark-results/m041-s04-t05-govulncheck.txt`
