---
id: T02
parent: S02
milestone: M038-pmw50e
key_files:
  - benchmark-results/fd-benchmark-m038-go-onnx-target-runtime.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T10:45:25.519Z
blocker_discovered: false
---

# T02: Ran the Go target-runtime performance benchmark against the actual Go ONNX API.

**Ran the Go target-runtime performance benchmark against the actual Go ONNX API.**

## What Happened

Started the Go ONNX API with namespace `m038-go-onnx-benchmark` and ran `benchmark.py` against `http://localhost:18000`. The benchmark completed with best cold latency `18.4ms`, warm latency mean `2.21ms`, max throughput `~895 req/s`, Batch L1 p95 `8.82ms`, and chunk reuse warm p95 `7.90ms`. The Redis L2 restart subcheck was skipped because restart command was intentionally empty for a bg_shell-managed server. Stopped the server and confirmed port 18000 is clean.

## Verification

Benchmark completed and artifact checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `BENCHMARK_API_URL=http://localhost:18000 ... uv run --python 3.13 --with requests --with redis python benchmark.py > benchmark-results/fd-benchmark-m038-go-onnx-target-runtime.txt` | 0 | ✅ pass — benchmark completed against Go ONNX endpoint | 25300ms |
| 2 | `gsd_exec M038 benchmark artifact summary/check` | 0 | ✅ pass — required config markers present, no leak markers/signed URLs | 52ms |
| 3 | `stop server and port check` | 0 | ✅ pass — no background processes, port_18000_clean | 0ms |

## Deviations

Set `BENCHMARK_API_RESTART_COMMAND=''` because the API was managed by bg_shell; Redis L2 restart subchecks were truthfully skipped and recorded in the benchmark summary.

## Known Issues

Benchmark.py calls Redis FLUSHALL; this side effect occurred during the benchmark. Redis L2 restart proof is skipped for this run because bg_shell-managed API restart was not wired into benchmark.py.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m038-go-onnx-target-runtime.txt`
