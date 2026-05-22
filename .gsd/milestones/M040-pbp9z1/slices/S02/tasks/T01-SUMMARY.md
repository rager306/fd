---
id: T01
parent: S02
milestone: M040-pbp9z1
key_files:
  - tools/run_m040_s02_docker_restart_proof.sh
  - tools/verify_m040_s02_artifacts.py
key_decisions:
  - Use a fixed S02 cache namespace `m040-s02-onnx-restart` and fixed proof container/network names so the restart command can restart only the API container while Redis remains alive.
  - Expose API and Redis only on localhost host bindings for the local proof path.
  - Keep artifact verification focused on semantics and redaction boundaries rather than latency thresholds.
duration: 
verification_result: passed
completed_at: 2026-05-22T07:43:05.530Z
blocker_discovered: false
---

# T01: Added a reproducible Docker ONNX restart proof runner and artifact verifier for M040 S02.

**Added a reproducible Docker ONNX restart proof runner and artifact verifier for M040 S02.**

## What Happened

Created `tools/run_m040_s02_docker_restart_proof.sh` to build/reuse the packaged ONNX API image as `fd-api:onnx-1024`, start Redis and the API in a fixed S02 Docker network with localhost-only host bindings, record preflight health/smoke/container evidence, run `benchmark.py` against `http://127.0.0.1:18000`, and configure `BENCHMARK_API_RESTART_COMMAND` as `docker restart fd-m040-s02-onnx-api` so Redis remains alive during API restarts. Created `tools/verify_m040_s02_artifacts.py` to validate the preflight and benchmark artifacts for ONNX runtime metadata, cache namespace isolation, configured restart command, non-skipped Section 5 and Section 6 Redis L2 restart evidence, safe sanitized config fields, smoke response shape, localhost port binding, and prohibited leak patterns. The wrapper reports clear `BLOCKER:` lines for missing Docker, occupied proof ports, unexpected health metadata, and readiness failures rather than producing false pass evidence.

## Verification

Ran the task-required shell syntax and Python compile checks. Also ran an inline fixture test that proves the verifier passes valid synthetic artifacts and rejects a benchmark artifact that reports skipped Batch L2 proof. Re-ran syntax/compile after the stale-container handling patch.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `bash -n tools/run_m040_s02_docker_restart_proof.sh && python3 -m py_compile tools/verify_m040_s02_artifacts.py` | 0 | ✅ pass | 59ms |
| 2 | `bash -n tools/run_m040_s02_docker_restart_proof.sh && python3 -m py_compile tools/verify_m040_s02_artifacts.py && inline verifier positive/negative fixture test` | 0 | ✅ pass | 155ms |
| 3 | `bash -n tools/run_m040_s02_docker_restart_proof.sh && python3 -m py_compile tools/verify_m040_s02_artifacts.py` | 0 | ✅ pass | 57ms |

## Deviations

Did not execute the full Docker benchmark proof in this task; the T01 done condition only required syntax/compile checks and verifier semantics, while the runner is designed for the later measured proof step.

## Known Issues

None.

## Files Created/Modified

- `tools/run_m040_s02_docker_restart_proof.sh`
- `tools/verify_m040_s02_artifacts.py`
