---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Author and link the runtime recommendation operating contract

Why: the user needs a fresh-operator-readable final recommendation that integrates S01's same-host HTTP contract, S02's packaged ONNX restart/cache/legal proof, and S03's bounded candidate defer result. Expected executor skills: api-design, design-an-interface, write-docs, verify-before-complete. Do: create `benchmark-results/fd-runtime-recommendation-m040-s04.md` with sections for Recommendation, Evidence Envelope, Operating Contract, Caveats and Required Operator Checks, Non-Goals, Source Artifacts, and Redaction. State the exact stance: use `deepvk/USER-bge-m3`; prefer packaged ONNX for same-host performance only when explicitly deployed with S01 `/health.runtime` metadata, artifact/tokenizer verification, optional `ONNX_RUNTIME_SHA256` runtime-library integrity, isolated `EMBEDDING_CACHE_VERSION`, and smoke `POST /v1/embeddings`; keep TEI as current/default until explicit switch; no request-level fallback; alternative models remain `defer_candidate`. Include S02 evidence values as evidence, not brittle thresholds: ONNX backend/model/dimensions/cache namespace, API-only restart with Redis L2 reuse, legal PASS/no-regression parity, cleanup audit, and the caveat that `runtime_library_verified=false` in S02 because `ONNX_RUNTIME_SHA256` was not set. Add a concise discoverability link from `docs/same-host-embedding-service-contract.md` to the final artifact without duplicating the contract. Failure Modes (Q5): if source artifacts are unavailable or conflict, do not invent replacement evidence; document the blocker or defer language and update the verifier expectation only if the source evidence truly changed. Load Profile (Q6): documentation only; no runtime load. Negative Tests (Q7): the artifact must be rejected by the verifier if it omits cache namespace isolation, treats hosted CI as required, accepts BAAI/E5 as replacements, or claims `/health` alone proves live inference readiness. Done when the full verifier passes against the final artifact and all S02/S03 inputs.

## Inputs

- `tools/verify_m040_s04_recommendation.py`
- `docs/same-host-embedding-service-contract.md`
- `benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-onnx-docker-preflight.txt`
- `benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt`
- `benchmark-results/fd-m040-s02-proof-audit.txt`
- `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`

## Expected Output

- `benchmark-results/fd-runtime-recommendation-m040-s04.md`
- `docs/same-host-embedding-service-contract.md`

## Verification

python3 tools/verify_m040_s04_recommendation.py --artifact benchmark-results/fd-runtime-recommendation-m040-s04.md --s02-benchmark benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt --s02-preflight benchmark-results/fd-m040-s02-onnx-docker-preflight.txt --s02-legal benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt --s02-audit benchmark-results/fd-m040-s02-proof-audit.txt --s03-gate benchmark-results/fd-legal-model-quick-gate-m040-s03.md

## Observability Impact

Publishes the decision inputs and operating caveats in a structured artifact that future operators and agents can inspect without rerunning runtime experiments.
