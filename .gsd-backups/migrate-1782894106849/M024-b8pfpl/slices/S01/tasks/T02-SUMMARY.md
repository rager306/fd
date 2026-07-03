---
id: T02
parent: S01
milestone: M024-b8pfpl
key_files:
  - benchmark-results/fd-benchmark-m024-onnx-docker1024.txt
key_decisions:
  - Use packaged runtime label `packaged-onnx1024-docker`.
  - Use benchmark cache namespace `m024-onnx-docker-benchmark`.
  - Use `docker restart fd-onnx-m024-bench >/dev/null` for Redis L2 restart measurement.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:26:54.477Z
blocker_discovered: false
---

# T02: Ran the packaged ONNX Docker performance benchmark successfully.

**Ran the packaged ONNX Docker performance benchmark successfully.**

## What Happened

Ran `benchmark.py` against the packaged ONNX Docker container on port 18000 via `uv run --python 3.13`. The benchmark completed successfully and wrote `benchmark-results/fd-benchmark-m024-onnx-docker1024.txt`. The artifact records runtime label `packaged-onnx1024-docker`, build tags `onnx,hf_tokenizers`, manifests, `ONNX_MAX_SEQUENCE_LENGTH=1024`, cache namespace `m024-onnx-docker-benchmark`, and the container restart command. Key metrics: best cold latency 7.6ms, warm latency mean 2.03ms, max throughput about 864 req/s, Redis L2 restart 3.36ms, batch L1 p95 8.91ms, batch L2 p95 5.47ms, chunk reuse warm p95 5.24ms.

## Verification

Benchmark exited 0 and artifact exists with expected packaged ONNX config markers and summary metrics.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `BENCHMARK_API_URL=http://localhost:18000 ... uv run --python 3.13 --with requests --with redis python benchmark.py > benchmark-results/fd-benchmark-m024-onnx-docker1024.txt` | 0 | ✅ pass — benchmark completed | 36200ms |
| 2 | `read benchmark-results/fd-benchmark-m024-onnx-docker1024.txt` | 0 | ✅ pass — summary includes best cold 7.6ms, warm mean 2.03ms, max throughput ~864 req/s, Redis L2 restart 3.36ms | 0ms |

## Deviations

The benchmark config records `BENCHMARK_ONNX_RUNTIME_LIBRARY=/opt/onnxruntime/libonnxruntime.so.1.26.0`, which is the in-container runtime path. The host-side metadata collector cannot hash that path, so `runtime.onnx_runtime_library.exists=false` in the artifact. The actual image smoke and benchmark prove the library is present inside the container.

## Known Issues

`benchmark.py` metadata collection runs on host, so it cannot hash in-container `/opt/onnxruntime/libonnxruntime.so.1.26.0`. Future packaged benchmark metadata could include Docker image ID/digest and in-container file hash.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m024-onnx-docker1024.txt`
