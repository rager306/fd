---
id: T03
parent: S04
milestone: M003-xx4yc3
key_files:
  - benchmark-results/fd-benchmark-baseline-py313.txt
  - benchmark-results/fd-runtime-stats-logs.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T08:24:04.495Z
blocker_discovered: false
---

# T03: Correlated benchmark results with container stats and logs.

**Correlated benchmark results with container stats and logs.**

## What Happened

Captured Docker stats and logs after benchmark and summarized benchmark signals. Post-benchmark stats showed API low memory (~9.6MiB), Redis low memory (~6.1MiB), and TEI using ~1.7GiB. Benchmark baseline under Python 3.13 showed cold latency scaling with text size, warm cache mean around 2.00ms, repeated cached p95 2.17ms/p99 4.72ms, and maximum throughput around 742 req/s at concurrency 16. API logs during throughput are very noisy because every successful embedding logs at INFO. TEI logs include historical ONNX fallback warnings and backend max batch request limit of 4, suggesting TEI/backend batching constraints matter for future cold-path/batch optimization.

## Verification

`docker stats --no-stream` and service logs were captured to `benchmark-results/fd-runtime-stats-logs.txt`; benchmark summary was parsed from `benchmark-results/fd-benchmark-baseline-py313.txt`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker stats --no-stream fd_api fd_tei fd_redis; docker compose logs api/tei/redis highlights | tee benchmark-results/fd-runtime-stats-logs.txt` | 0 | ✅ pass: stats/log artifact captured | 0ms |
| 2 | `parse benchmark-results/fd-benchmark-baseline-py313.txt summary` | 0 | ✅ pass: warm mean 2.00ms, cached p95 2.17ms, max throughput ~742 req/s | 0ms |

## Deviations

The runtime stats/log capture produced a large log artifact; the persisted file `benchmark-results/fd-runtime-stats-logs.txt` contains the full evidence.

## Known Issues

Redis still reports host memory overcommit warning. TEI logs show historical ONNX-missing fallback to Candle CPU and `max_batch_requests=4`; service remains healthy.

## Files Created/Modified

- `benchmark-results/fd-benchmark-baseline-py313.txt`
- `benchmark-results/fd-runtime-stats-logs.txt`
