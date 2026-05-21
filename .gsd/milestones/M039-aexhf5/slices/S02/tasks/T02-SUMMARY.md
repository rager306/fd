---
id: T02
parent: S02
milestone: M039-aexhf5
key_files:
  - benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-21T11:29:18.301Z
blocker_discovered: false
---

# T02: Packaged Docker ONNX performance benchmark passed through actual packaged Go endpoint.

**Packaged Docker ONNX performance benchmark passed through actual packaged Go endpoint.**

## What Happened

Started packaged image `fd-api:onnx1024-m039` with `ONNX_RUNTIME_SHA256` and namespace `m039-docker-benchmark`. Health confirmed backend `onnx`, artifact/tokenizer/runtime library verified, CPU provider, dimensions 1024, and expected namespace. Ran benchmark.py against `http://localhost:18000`. Results: best cold latency `10.3ms`, warm latency mean `1.52ms`, max throughput `~937 req/s`, Batch L1 p95 `4.78ms`, chunk reuse warm p95 `13.55ms`, Redis L2 restart skipped. Artifact checks passed after using metric-based markers. Container was stopped and port 18000 is clean.

## Verification

Benchmark, artifact checks, and cleanup passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker run fd-api:onnx1024-m039 with namespace m039-docker-benchmark and ONNX_RUNTIME_SHA256` | 0 | ✅ pass — container started | 6100ms |
| 2 | `packaged benchmark health precheck` | 0 | ✅ pass — artifact/tokenizer/runtime verified, namespace m039-docker-benchmark | 4200ms |
| 3 | `BENCHMARK_API_URL=http://localhost:18000 ... uv run --python 3.13 --with requests --with redis python benchmark.py > benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt` | 0 | ✅ pass — benchmark completed; max throughput ~937 req/s | 23100ms |
| 4 | `strict artifact marker check` | 1 | ⚠️ fail — check expected absent `Benchmark completed` marker; metrics were present | 5000ms |
| 5 | `corrected benchmark artifact checks` | 0 | ✅ pass — required metric/config markers present, no leak markers/signed URLs | 3400ms |
| 6 | `docker rm -f fd-onnx-m039-benchmark && port check` | 0 | ✅ pass — port_18000_clean | 5000ms |

## Deviations

Set `BENCHMARK_API_RESTART_COMMAND=''` because the packaged container was managed externally for this benchmark; Redis L2 restart checks were truthfully skipped. An initial artifact check expected a non-existent `Benchmark completed` marker and failed, then the corrected metric-based artifact check passed without rerunning benchmark.

## Known Issues

Benchmark.py calls Redis FLUSHALL. The benchmark artifact's runtime snapshot reports default labels/configured=false for local ONNX paths because the packaged image supplies ONNX artifacts internally rather than through host path env vars; endpoint evidence and health precheck establish packaged ONNX runtime.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt`
