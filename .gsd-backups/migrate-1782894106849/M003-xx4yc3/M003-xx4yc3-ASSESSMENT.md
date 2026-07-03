# Runtime Validation and Performance Assessment

## Baseline Environment

- Stack: Docker Compose with `api`, `tei`, and `redis`.
- Benchmark runtime: `uv run --python 3.13 --with requests --with redis`.
- Python version: 3.13.12.
- Benchmark artifact: `benchmark-results/fd-benchmark-baseline-py313.txt`.
- Runtime stats/log artifact: `benchmark-results/fd-runtime-stats-logs.txt`.

## Runtime Findings

### Fixed during validation

1. **Redis host exposure in override**
   - Evidence: Redis logs contained external attack attempts while Redis was previously published on `0.0.0.0:6379`.
   - Fix: `docker-compose.override.yaml` now binds Redis to `127.0.0.1:6379:6379`.
   - Result: Local benchmark access remains, but Redis is not exposed on all host interfaces.

2. **Stale `fd_api` container name conflict**
   - Evidence: first `docker compose up -d --build` failed because an exited `fd_api` container existed.
   - Fix: removed stale container with `docker rm fd_api`.
   - Result: Compose stack starts; `fd_api`, `fd_redis`, and `fd_tei` healthy.

### Non-blocking operational notes

1. **Redis host memory overcommit warning**
   - Evidence: Redis logs warn `Memory overcommit must be enabled`.
   - Scope: host configuration, not app code.
   - Recommendation: document/set `vm.overcommit_memory=1` on deployment hosts if Redis persistence/background save matters.

2. **TEI ONNX fallback to Candle CPU**
   - Evidence: TEI logs show missing ONNX artifacts and fallback to Candle CPU backend.
   - Scope: model artifact/runtime optimization.
   - Recommendation: evaluate exporting ONNX artifacts only if cold-path TEI latency is a bottleneck worth optimizing.

3. **API success logs are noisy under throughput**
   - Evidence: API logs one `embeddings generated` INFO line per request; benchmark throughput generated large log volume.
   - Scope: observability/cost, not correctness.
   - Recommendation: keep error/warn logs eager, but sample or debug-gate success logs under load.

## Python 3.13 Benchmark Baseline

Accepted command:

```bash
uv run --python 3.13 --with requests --with redis python benchmark.py
```

Key results:

| Scenario | Result |
|---|---:|
| Best cold latency | 19.5ms |
| Warm latency mean | 2.00ms |
| Repeated cached mean | 1.55ms |
| Repeated cached p95 | 2.17ms |
| Repeated cached p99 | 4.72ms |
| Max measured throughput | ~742 req/s at 16 concurrent |

Cold/warm latency by text size:

| Text | Cold mean | Warm mean | Speedup |
|---|---:|---:|---:|
| short | 19.5ms | 1.62ms | 12.0x |
| medium | 48.7ms | 1.82ms | 26.8x |
| long | 118.2ms | 2.03ms | 58.3x |
| very_long | 132.7ms | 2.53ms | 52.4x |

Throughput table:

| Concurrency | Req/s | Mean | p50 | p95 |
|---:|---:|---:|---:|---:|
| 1 | 575.6 | 1.5ms | 1.4ms | 2.3ms |
| 4 | 830.6 | 4.1ms | 3.9ms | 6.6ms |
| 8 | 777.2 | 9.1ms | 8.6ms | 15.8ms |
| 16 | 742.0 | 19.4ms | 18.4ms | 33.4ms |

Note: benchmark summary says max throughput is `~742 req/s (16 concurrent)`, but the table shows `830.6 req/s` at concurrency 4. This is a benchmark-script reporting bug and should be fixed before using the summary line for automated comparisons.

## Prioritized Optimization Plan

### P1: Fix benchmark throughput summary selection

- Evidence: table max is 830.6 req/s at concurrency 4, while summary says 742 req/s at concurrency 16.
- Expected outcome: benchmark summary reliably reports true max throughput and associated concurrency.
- Likely file: `benchmark.py`.
- Validation: rerun `uv run --python 3.13 --with requests --with redis python benchmark.py` and compare summary to table max.
- Rollback: revert only benchmark.py summary logic.

### P2: Add cache metrics and log sampling

- Evidence: cache hits deliver ~1.5-2.0ms mean latency, but runtime logs do not expose structured L1/L2 hit counts or latencies; successful requests log per request at INFO.
- Expected outcome: operators can distinguish L1 hit, L2 hit, cold TEI call, and cache write failures without high-volume INFO spam.
- Likely files: `api/cache/tiered.go`, `api/handlers/embeddings.go`, `api/handlers/batch.go`.
- Validation: unit tests for counters/log path plus smoke requests showing cache hit/miss signals.
- Rollback: remove metrics/log sampling wrapper without changing cache semantics.

### P3: Add benchmark modes for cold path, warm path, and cache-disabled path

- Evidence: current benchmark uses Redis FLUSHALL and repeated hits, but does not isolate L1 vs L2 vs TEI cost or expose cache-disabled baseline.
- Expected outcome: optimization work can prove whether gains affect API overhead, Redis, or TEI.
- Likely file: `benchmark.py`; possibly API env/config if cache-disable switch is introduced.
- Validation: benchmark outputs separate sections for L1 warm, L2 warm after API restart, and cold TEI.
- Rollback: keep current simple benchmark as default; make deeper modes optional flags.

### P4: Evaluate TEI backend/model artifact optimization

- Evidence: TEI uses ~1.7GiB and logs ONNX-missing fallback; cold path has 19.5-132.7ms means for tested texts, with historical TEI timings showing higher latency under some loads.
- Expected outcome: decide whether ONNX export/runtime tuning improves cold-path latency enough to justify model artifact complexity.
- Likely files/config: TEI image/compose args, model volume contents, docs.
- Validation: A/B benchmark against same Python 3.13 benchmark and same text set; compare cold p50/p95 and memory.
- Rollback: revert TEI config/model artifact change and keep current Candle backend.

### P5: Tune handler batching only after cache metrics exist

- Evidence: TEI backend enforces `max_batch_requests=4`; current `/embeddings/batch` is functional, but future performance work should avoid increasing API-side batch fan-out blindly.
- Expected outcome: batch throughput improves without queue-time regressions or worse p95.
- Likely files: `api/handlers/batch.go`, `api/embed/tei.go`.
- Validation: batch benchmark with varied batch sizes and dimensions; confirm response shape and p95.
- Rollback: restore current sequential/simple batch behavior.

## Recommended Next Milestone

Create a focused optimization milestone rather than mixing more changes into runtime validation:

1. Fix `benchmark.py` max-throughput summary bug.
2. Add cache hit/miss metrics/log sampling.
3. Extend benchmark modes to separate L1, L2, and cold TEI.
4. Only then tune TEI/backend batching based on new evidence.
