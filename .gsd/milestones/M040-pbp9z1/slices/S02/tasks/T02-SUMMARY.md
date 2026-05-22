---
id: T02
parent: S02
milestone: M040-pbp9z1
key_files:
  - tools/run_m040_s02_docker_restart_proof.sh
  - tools/verify_m040_s02_artifacts.py
  - benchmark-results/fd-m040-s02-onnx-docker-preflight.txt
  - benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt
key_decisions:
  - Validate smoke readiness by response model plus 1024 dimensions, accepting the API's numeric embedding vector shape instead of requiring base64 encoding.
  - Use a local PYTHONPATH package target for the benchmark redis dependency rather than overriding PEP 668-managed system Python.
duration: 
verification_result: passed
completed_at: 2026-05-22T07:53:03.766Z
blocker_discovered: false
---

# T02: Ran the packaged ONNX Docker restart proof and produced passing Redis L2 restart benchmark/preflight artifacts.

**Ran the packaged ONNX Docker restart proof and produced passing Redis L2 restart benchmark/preflight artifacts.**

## What Happened

Executed the packaged ONNX Docker restart/cache proof for S02. Initial runs exposed local proof-runner and environment issues before the benchmark could complete: the runner did not print successful /health JSON from wait_for_http_health, json_get consumed stdin via a heredoc instead of the piped JSON body, and the smoke assertion expected base64 even though the API returned a numeric 1024-dimensional embedding vector. I patched those narrow proof-runner/verifier mismatches, installed the missing Python redis dependency into a local PYTHONPATH target without modifying system Python, and reran the wrapper end-to-end. The final wrapper built/reused `fd-api:onnx-1024`, reused the proof Redis container without restarting it, started the ONNX API on `http://127.0.0.1:18000`, recorded /health runtime identity and smoke embedding shape, ran `benchmark.py` with `BENCHMARK_API_RESTART_COMMAND='docker restart fd-m040-s02-onnx-api'`, measured Redis L2 reuse after API-only restart, and completed artifact verification successfully. Final benchmark evidence includes `Redis L2 restart: 4.73ms after API restart`, `Batch L2 p95: 3.89ms after API restart`, `redis_delta/l2_after_api_restart: hits=1 misses=0 evicted=0 expired=0 commands=3 db0_key_delta=0`, and `redis_delta/batch_l2_after_api_restart: hits=16 misses=0 evicted=0 expired=0 commands=18 db0_key_delta=0`.

## Verification

Verified shell/Python syntax after runner/verifier fixes, installed and confirmed the local redis Python dependency for benchmark.py, ran the full Docker proof wrapper successfully with local PYTHONPATH, reran the task-required artifact verifier command, checked both expected artifacts are non-empty, and confirmed the API proof container was cleaned up while Redis remained alive for the API-only restart proof.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `bash -n tools/run_m040_s02_docker_restart_proof.sh && python3 -m py_compile tools/verify_m040_s02_artifacts.py` | 0 | ✅ pass | 60ms |
| 2 | `python3 -m pip install --target .gsd/runtime/python-packages redis && PYTHONPATH=.gsd/runtime/python-packages python3 -c 'import redis; print(redis.__version__)'` | 0 | ✅ pass | 13670ms |
| 3 | `PYTHONPATH=.gsd/runtime/python-packages ./tools/run_m040_s02_docker_restart_proof.sh` | 0 | ✅ pass | 50087ms |
| 4 | `python3 tools/verify_m040_s02_artifacts.py --benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt --preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt` | 0 | ✅ pass | 69ms |
| 5 | `test -s benchmark-results/fd-m040-s02-onnx-docker-preflight.txt && test -s benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt && docker ps checks for proof containers` | 0 | ✅ pass | 122ms |

## Deviations

Had to apply narrow fixes to the T01 runner/verifier before the proof could complete: shell helpers now emit captured /health JSON and parse piped JSON correctly; smoke checks now validate the observable 1024-dimension embedding shape rather than requiring base64; verifier now reads benchmark.py's sanitized environment from environment.values. The local system Python lacked redis and python3-venv, so redis was installed into .gsd/runtime/python-packages and the final proof was run with PYTHONPATH pointing there.

## Known Issues

The proof Redis container `fd-m040-s02-redis` remains running as intended for the API-only restart proof; the API container is cleaned up by the runner trap. Cleanup evidence is left for the slice cleanup task if applicable.

## Files Created/Modified

- `tools/run_m040_s02_docker_restart_proof.sh`
- `tools/verify_m040_s02_artifacts.py`
- `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt`
- `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt`
