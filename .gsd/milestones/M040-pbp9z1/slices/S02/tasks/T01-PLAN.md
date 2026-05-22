---
estimated_steps: 13
estimated_files: 2
skills_used: []
---

# T01: Create reproducible S02 Docker proof runner and artifact verifier

skills_used:
  - api-design
  - observability
  - tdd
  - verify-before-complete

Why: S02 should not depend on ad-hoc shell history. A small proof runner and verifier make the Docker restart/cache proof repeatable while preserving the existing `benchmark.py` benchmark semantics.

Threat Surface (Q3): The runner will execute Docker commands and a benchmark restart hook; keep commands hard-coded or allowlisted for the local proof container only, do not interpolate untrusted input, do not print secrets, bind the API to localhost proof port 18000, and keep Redis on localhost.

Requirement Impact (Q4): Advances R002 and R004 directly; supports R001 and R003 by verifying runtime identity and safe config. Re-verify benchmark artifact semantics, health metadata, smoke inference, and redaction boundaries.

Failure Modes (Q5): Missing Docker permission, missing ONNX artifact/runtime/tokenizer, missing Redis, occupied port 18000, failed image build, failed `/health`, failed smoke embedding, failed restart hook, or skipped benchmark sections must produce clear blocker text in the proof artifact rather than a false pass.

Load Profile (Q6): The runner is a one-host benchmark path. Shared resources are Docker daemon, Redis L2 cache, port 18000, and CPU/RAM for ONNX inference. At 10x concurrent runs, port/container/cache namespace collisions break first; enforce a fixed S02 namespace and single proof container.

Negative Tests (Q7): Verifier must fail when benchmark artifact is missing, when `api_restart_command_configured` is false, when Section 5 or Section 6 is skipped, when sanitized config has `input_texts_logged` other than false, or when audit artifacts contain known prohibited secret/text patterns.

Do: Add a lightweight `tools/run_m040_s02_docker_restart_proof.sh` wrapper that builds/reuses `fd-api:onnx-1024`, starts the packaged ONNX API on `127.0.0.1:18000` with `EMBEDDING_CACHE_VERSION=m040-s02-onnx-restart`, records preflight health/smoke/container evidence to `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt`, runs `benchmark.py` with a restart command that restarts only the API container while leaving Redis alive, and writes `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt`. Add `tools/verify_m040_s02_artifacts.py` to assert required artifact semantics and leak/audit rules. Do not add first-class Docker lifecycle logic to `benchmark.py` unless the wrapper proves insufficient.

Done when: Shell syntax and Python compile checks pass, and the verifier has explicit checks for configured restart command, non-skipped Redis L2 restart sections, non-skipped batch L2 restart sections, safe config, expected runtime/cache fields, and prohibited leak patterns.

## Inputs

- `benchmark.py`
- `tools/build_onnx_image.sh`
- `Dockerfile.onnx`
- `docs/same-host-embedding-service-contract.md`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt`
- `README.md`

## Expected Output

- `tools/run_m040_s02_docker_restart_proof.sh`
- `tools/verify_m040_s02_artifacts.py`

## Verification

bash -n tools/run_m040_s02_docker_restart_proof.sh && python3 -m py_compile tools/verify_m040_s02_artifacts.py

## Observability Impact

Adds local proof artifacts that expose phase, command, health metadata, smoke result, benchmark semantics, and explicit failure reasons without secrets.
