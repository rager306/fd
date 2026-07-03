---
estimated_steps: 11
estimated_files: 2
skills_used: []
---

# T02: Run packaged ONNX Docker restart and Redis L2 benchmark proof

skills_used:
  - observability
  - verify-before-complete

Why: The core S02 risk is that M039 skipped Redis L2 restart proof because `BENCHMARK_API_RESTART_COMMAND` was empty. This task runs the actual packaged ONNX API and proves Redis L2 survives API-only restart.

Threat Surface (Q3): Real Docker and Redis are used. Redis must survive the API restart and must not be exposed beyond localhost. The benchmark calls Redis `FLUSHALL`; run only against the local proof Redis, never shared/production Redis.

Requirement Impact (Q4): Primary proof for R002 and R004; supports R003 by exercising environment-configured runtime/cache parameters. Re-verify S01 readiness contract via `/health` plus smoke embedding before accepting benchmark evidence.

Failure Modes (Q5): If Docker, image build, Redis, ONNX runtime verification, health, smoke inference, restart command, or benchmark sections fail, record exact command/output in the preflight or benchmark artifact and stop. A skipped Section 5 or Section 6 is a task failure, not success.

Load Profile (Q6): Benchmark intentionally stresses cache paths, batch paths, and API restart. Shared resources are Redis, Docker, port 18000, and CPU. Cache namespace isolation prevents TEI/ONNX contamination.

Negative Tests (Q7): Verifier must reject artifacts where the restart command is empty, Redis L2 restart is skipped, batch L2 restart is skipped, runtime backend is not ONNX, smoke dimensions are not 1024, cache namespace is not `m040-s02-onnx-restart`, or raw input texts/secrets appear.

Do: Execute the runner from T01. Confirm preflight shows packaged ONNX API started on `http://localhost:18000`, `/health.runtime.backend=onnx`, `/health.runtime.model=deepvk/USER-bge-m3`, `/health.runtime.dimensions=1024`, expected cache namespace, and `runtime_library_verified=true` when `ONNX_RUNTIME_SHA256` is present. Confirm smoke `/v1/embeddings` returns response model `deepvk/USER-bge-m3` and 1024 dimensions. Run `benchmark.py` through the wrapper with isolated namespace and API-only Docker restart command. Do not restart Redis.

Done when: `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt` and `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt` exist, and the verifier passes on both artifacts with measured single-request and batch Redis L2 restart sections.

## Inputs

- `tools/run_m040_s02_docker_restart_proof.sh`
- `tools/verify_m040_s02_artifacts.py`
- `benchmark.py`
- `tools/build_onnx_image.sh`
- `Dockerfile.onnx`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `docs/same-host-embedding-service-contract.md`

## Expected Output

- `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt`

## Verification

python3 tools/verify_m040_s02_artifacts.py --benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt --preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt

## Observability Impact

Produces the slice's main operational evidence: runtime identity, smoke inference readiness, Docker restart command configuration, Redis L2 restart latency, batch L2 restart latency, and sanitized benchmark config.
