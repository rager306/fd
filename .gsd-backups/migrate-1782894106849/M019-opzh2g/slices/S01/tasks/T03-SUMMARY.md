---
id: T03
parent: S01
milestone: M019-opzh2g
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m019-onnx1024.txt
key_decisions:
  - M019 benchmark artifact must include `ONNX_MAX_SEQUENCE_LENGTH=1024` for comparability.
  - ONNX 1024 benchmark result: best cold latency 8.3ms, warm latency mean 1.19ms, max throughput about 858 req/s, Redis L2 restart 2.02ms, batch L1 p95 4.09ms, batch L2 p95 6.70ms, chunk reuse warm p95 6.51ms.
duration: 
verification_result: passed
completed_at: 2026-05-20T08:16:14.223Z
blocker_discovered: false
---

# T03: Ran the ONNX 1024 benchmark and captured sanitized performance metrics.

**Ran the ONNX 1024 benchmark and captured sanitized performance metrics.**

## What Happened

Ran the ONNX 1024 benchmark with isolated namespace `m019-onnx-1024-benchmark`, actual ONNX restart command, and sanitized config snapshot. After improving benchmark metadata allowlist, the final artifact records `ONNX_MAX_SEQUENCE_LENGTH=1024` and ONNX runtime settings. Artifact hygiene check confirms no raw benchmark text leaks.

## Verification

Benchmark exited 0, summary metrics were extracted, config includes ONNX sequence length, and artifact hygiene passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `BENCHMARK_API_URL=http://localhost:18000 ... ONNX_MAX_SEQUENCE_LENGTH=1024 ... uv run --python 3.13 --with requests --with redis python benchmark.py > benchmark-results/fd-benchmark-m019-onnx1024.txt` | 0 | ✅ pass — benchmark artifact written | 35000ms |
| 2 | `python benchmark artifact summary/hygiene checks` | 0 | ✅ pass — ONNX_MAX_SEQUENCE_LENGTH recorded; raw_benchmark_text_leaks=0 | 0ms |

## Deviations

During first benchmark run, `benchmark.py` did not include `ONNX_MAX_SEQUENCE_LENGTH` in its sanitized env allowlist, so I updated the allowlist and reran the benchmark. I also hardened the runtime helper to kill the actual `fd-api` listener left by `go run` child processes before restarts.

## Known Issues

Benchmark validates performance on current local KVM/QEMU host only. It does not validate Docker/CI packaging or production rollout.

## Files Created/Modified

- `benchmark.py`
- `benchmark-results/fd-benchmark-m019-onnx1024.txt`
