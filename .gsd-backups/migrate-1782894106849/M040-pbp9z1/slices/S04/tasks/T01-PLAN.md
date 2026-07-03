---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Add semantic verifier for the final recommendation artifact

Why: S04's highest-risk failure is a plausible narrative that violates the evidence envelope by recommending an unproven candidate model, treating `/health` as live inference readiness, omitting Redis cache namespace isolation, overclaiming legal quality, leaking raw text/secrets, or making hosted CI a readiness gate. Expected executor skills: api-design, design-an-interface, tdd, verify-before-complete. Do: create `tools/verify_m040_s04_recommendation.py` using the existing S02/S03 verifier style; implement required-heading checks, source-artifact reference checks, explicit recommendation stance checks, S02 packaged ONNX legal/restart/cache/preflight semantic checks, S03 `defer_candidate` checks, redaction/prohibited-pattern checks, and a `--self-test` that exercises at least one passing fixture plus negative fixtures for missing cache isolation, accepted candidate replacement, and hosted-CI-as-gate language. Failure Modes (Q5): if an evidence file is missing, fail with a named path and non-zero exit; if an evidence file is malformed or lacks required markers, fail closed rather than warning; if arguments are incomplete, show usage and exit non-zero. Load Profile (Q6): small local text scans only; per-operation cost is reading a handful of benchmark/docs files; 10x artifact size should still be bounded by straightforward line/string scanning, with no network/runtime calls. Negative Tests (Q7): self-test must include malformed/missing required sections, unproven alternative-model acceptance, hosted CI gate language, and prohibited secret/raw-text markers. Done when the verifier compiles and its self-test passes without needing the final artifact to exist yet.

## Inputs

- `tools/verify_m040_s02_artifacts.py`
- `tools/verify_legal_model_quick_gate_artifact.py`
- `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt`
- `benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-proof-audit.txt`
- `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`
- `docs/same-host-embedding-service-contract.md`

## Expected Output

- `tools/verify_m040_s04_recommendation.py`

## Verification

python3 -m py_compile tools/verify_m040_s04_recommendation.py
python3 tools/verify_m040_s04_recommendation.py --self-test

## Observability Impact

Creates the machine-checkable diagnostic surface for S04: failures should name the violated section/evidence/caveat so later agents can repair artifact drift quickly.
