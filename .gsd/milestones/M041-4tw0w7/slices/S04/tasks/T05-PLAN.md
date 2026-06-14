---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T05: Validated S04 performance under the accepted D045 cache-hot steady-state contract, with real cache-miss diagnostics preserved.

Run tools/verify_fd_v2_perf.sh against a current fd instance backed by real inference. The verifier must prewarm each measured payload through real inference, then validate cache-hot T-P-1..T-P-5 latency/error targets with `X-Cache: HIT`. It must also include non-blocking cache-miss diagnostics so TEI CPU miss latency remains visible without blocking S04.

## Inputs

- `tools/verify_fd_v2_perf.sh`
- `docs/fd-v2.md`
- `benchmark-results/m041-s04-t04-bottleneck-measurement.txt`
- `benchmark-results/m041-s04-t05-cache-namespace.txt`

## Expected Output

- `benchmark-results/fd-v2-perf-validation-m041-s04.md`

## Verification

FD_BASE_URL=http://localhost:8000 tools/verify_fd_v2_perf.sh exits 0 against current fd. Artifact contains p50/p95/p99 for cache-hot T-P cases, 100 sequential 0 errors, 4x8 concurrent <2s, cache HIT <5ms, and non-blocking cache-miss diagnostics.
