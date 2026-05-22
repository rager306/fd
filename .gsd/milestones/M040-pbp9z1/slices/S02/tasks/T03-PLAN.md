---
estimated_steps: 11
estimated_files: 2
skills_used: []
---

# T03: Run legal guard, leak audit, and cleanup proof

skills_used:
  - observability
  - verify-before-complete

Why: S02 also owns R001 and R004. Restart/cache speed is not acceptable unless the packaged ONNX runtime remains legal-domain compatible and artifacts are safe to hand to S04.

Threat Surface (Q3): Legal evaluator sends corpus-derived text to local APIs and writes result artifacts. Artifacts must not include raw legal/probe texts or secrets. Cleanup must not remove unrelated containers or Redis state outside the proof scope.

Requirement Impact (Q4): Validates or truthfully blocks R001 for the S02 packaged proof; closes R004 artifact-safety evidence; supports R006 for final S04 recommendation.

Failure Modes (Q5): If TEI baseline API is unavailable, ONNX API has already stopped, legal evaluator fails thresholds, or cleanup cannot remove the proof container, write exact blocker evidence to the audit artifact. If legal rerun is impossible, the audit must explicitly cite whether M039 legal PASS is still reusable by matching runtime/image/artifact inputs, not silently assume it.

Load Profile (Q6): Legal evaluation performs batched embedding requests over the structured test corpus. Shared resources are TEI API, ONNX API, Redis, CPU, and cache namespace. At higher load, API availability and Redis memory pressure break first; keep this as a bounded single-run gate.

Negative Tests (Q7): Final verifier must fail when legal result is missing or non-PASS, benchmark/preflight files are missing, leak audit finds prohibited patterns, cleanup evidence is absent, proof container remains running, or port 18000 remains bound.

Do: Run `tools/evaluate_legal_retrieval.py` against the baseline TEI API and packaged ONNX API with `--onnx-api-url http://localhost:18000`, `--onnx-runtime-label docker-onnx-go-api-m040-s02`, and `--onnx-cache-namespace m040-s02-onnx-restart`, saving `benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt`. Run the verifier over benchmark, preflight, and legal artifacts. Record `benchmark-results/fd-m040-s02-proof-audit.txt` with leak-audit verdict, cleanup commands/output, Docker container status, port 18000 clear/blocked status, and any blocker details. Stop/remove only the S02 proof ONNX API container; leave Redis in the state needed for local development unless the runner explicitly created an isolated Redis.

Done when: legal result PASS or a truthful blocker is recorded, leak audit passes, proof container is gone, port 18000 is clear, and the verifier passes over all S02 artifacts.

## Inputs

- `tools/verify_m040_s02_artifacts.py`
- `tools/run_m040_s02_docker_restart_proof.sh`
- `tools/evaluate_legal_retrieval.py`
- `tests/44-FZ-2026-articles.jsonl`
- `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt`
- `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt`

## Expected Output

- `benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-proof-audit.txt`

## Verification

python3 tools/verify_m040_s02_artifacts.py --benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt --preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt --legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt --audit benchmark-results/fd-m040-s02-proof-audit.txt

## Observability Impact

Adds final closeout surfaces for legal no-regression, artifact redaction, blocker truthfulness, and host cleanup state.
