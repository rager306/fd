---
id: T01
parent: S01
milestone: M024-b8pfpl
key_files: []
key_decisions:
  - Use packaged image `fd-api:onnx1024-m022-final` with container name `fd-onnx-m024-bench`.
  - Use isolated cache namespace `m024-onnx-docker-benchmark`.
  - Use `docker restart fd-onnx-m024-bench` as `BENCHMARK_API_RESTART_COMMAND`; a smoke restart succeeded after brief expected connection refusals during restart.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:25:17.615Z
blocker_discovered: false
---

# T01: Prepared the packaged ONNX Docker benchmark target and restart command.

**Prepared the packaged ONNX Docker benchmark target and restart command.**

## What Happened

Prepared the packaged ONNX benchmark target. The packaged ONNX image is running on port 18000, `/health` passed, a smoke embedding returned 1024 dimensions for `deepvk/USER-bge-m3`, and the restart command for benchmark L2 checks successfully restarted the same container and returned to health.

## Verification

Health, smoke embedding, and restart command smoke passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker run --rm --name fd-onnx-m024-bench --network host ... fd-api:onnx1024-m022-final` | 0 | ✅ pass — ready on port 18000 | 7000ms |
| 2 | `curl -fsS http://localhost:18000/health` | 0 | ✅ pass — onnx_docker_18000_health=pass | 0ms |
| 3 | `curl -fsS http://localhost:18000/v1/embeddings ...` | 0 | ✅ pass — embedding_dims=1024 model=deepvk/USER-bge-m3 | 0ms |
| 4 | `docker restart fd-onnx-m024-bench && curl --retry ... /health` | 0 | ✅ pass — restart_command_smoke=pass | 0ms |

## Deviations

None.

## Known Issues

`benchmark.py` will call Redis `FLUSHALL`; this is accepted benchmark behavior and affects the shared Redis instance.

## Files Created/Modified

None.
