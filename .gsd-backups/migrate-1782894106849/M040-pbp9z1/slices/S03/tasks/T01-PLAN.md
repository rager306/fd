---
estimated_steps: 16
estimated_files: 2
skills_used: []
---

# T01: Add bounded legal model gate tooling

Expected executor skills for task-plan frontmatter: grill-me, sentence-transformers, api-design, write-docs, verify-before-complete.

Why: The existing `tools/evaluate_legal_retrieval.py` is a same-model TEI-vs-ONNX parity evaluator; S03 needs a replacement-candidate gate that compares retrieval metrics and availability without misusing cross-model cosine or top-1 parity. This task creates a small dedicated tool/verifier so T02 can produce a bounded, redacted artifact whether candidates run successfully or must be truthfully deferred.

Do:
1. Add `tools/evaluate_legal_model_quick_gate.py` as a CLI wrapper around the existing legal corpus/evaluator patterns, without modifying the service API and without relying on the `/v1/embeddings` request `model` field for selection.
2. Reuse safe ideas from `tools/evaluate_legal_retrieval.py`: corpus hash, deterministic document/query selection, text hashing/redaction, dimension checks, batch embedding requests, recall@k/MRR metrics, and markdown artifact rendering.
3. Support a dry-run/availability-only path that records candidate shortlist and stop reasons without live embedding calls; support at most two candidates and fail fast if more are provided.
4. For live runs, require separate baseline/candidate endpoint configuration, model IDs, expected dimensions, runtime labels, and cache namespaces; smoke `/health` and one `/v1/embeddings` request before legal corpus calls.
5. Use retrieval metrics and operational compatibility to render one final outcome (`keep_current`, `reject_candidate`, or `defer_candidate`); explicitly mark cross-model cosine/parity as not applicable for different models.
6. Add `tools/verify_legal_model_quick_gate_artifact.py` to validate the rendered artifact shape, candidate count, required metadata, legal redaction statement, absence of raw-text sections, absence of obvious secret/token patterns, required verdict, and stop-reason/metrics consistency. Include a `--self-test` mode that builds temporary sample artifacts and exercises pass/fail cases without reading `.gsd` or gitignored paths.
7. If modifying any existing function/class/method instead of only adding new files, first run GitNexus impact analysis for the modified symbol and preserve all direct callers; otherwise keep the change additive.

Done when: The new tools compile, self-test, and can render/validate a dry-run artifact with no raw corpus text, no candidate count above two, and no cross-model cosine acceptance claim.

Threat Surface (Q3): candidate endpoints are local HTTP services and model IDs are operator input; malformed URLs, wrong dimensions, timeouts, and invalid JSON must fail closed into a defer/reject stop reason rather than producing a misleading PASS.
Requirement Impact (Q4): owns R008, supports R001 and R006, preserves D040 and D041. Re-verify the legal artifact redaction shape, bounded candidate count, deployment-scoped model semantics, and no API contract changes.
Failure Modes (Q5): endpoint refused/timeout -> candidate deferred with exact phase; malformed health/embedding response -> candidate rejected/deferred with phase and sanitized error; dimension mismatch -> reject/defer candidate and stop legal metric comparison; corpus parse failure -> hard fail.
Load Profile (Q6): bounded corpus run should cap documents and title/self queries; batch size and timeout should be explicit; at 10x load, HTTP runtime and memory would break first, which is out of scope for quick gate.
Negative Tests (Q7): self-test must cover too many candidates, missing required metadata, raw_text_logged not false, artifact containing prohibited secret-like patterns, missing verdict, and cross-model cosine being treated as an acceptance metric.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `tests/44-FZ-2026-articles.jsonl`
- `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt`
- `docs/same-host-embedding-service-contract.md`
- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `api/main.go`

## Expected Output

- `tools/evaluate_legal_model_quick_gate.py`
- `tools/verify_legal_model_quick_gate_artifact.py`
- `benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md`

## Verification

python3 -m py_compile tools/evaluate_legal_model_quick_gate.py tools/verify_legal_model_quick_gate_artifact.py
python3 tools/verify_legal_model_quick_gate_artifact.py --self-test
python3 tools/evaluate_legal_model_quick_gate.py --dry-run --corpus tests/44-FZ-2026-articles.jsonl --output benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md --baseline-model deepvk/USER-bge-m3 --baseline-runtime-label tei-default --baseline-dimensions 1024 --candidate-model BAAI/bge-m3 --candidate-runtime-label candidate-bge-m3 --candidate-dimensions 1024 --candidate-cache-namespace m040-s03-candidate-bge-m3 --max-docs 32 --max-title-queries 8 --max-self-queries 8
python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md --max-candidates 2

## Observability Impact

Adds file-based diagnostics for model-gate runs: selected corpus counts, model IDs, dimensions, runtime labels, cache namespaces, health/smoke status, phase-specific stop reasons, metrics, verdict, and redaction status. No production API observability changes.
