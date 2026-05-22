# S03: Bounded legal model quick gate

**Goal:** Bound alternative embedding-model exploration to at most two plausible Russian/legal candidates, compare them against the current deepvk/USER-bge-m3 same-host service using legal-domain evidence or a truthful availability deferral, and produce a sanitized quick-gate artifact for S04 without changing the service API contract.
**Demo:** After this, alternative model scope is bounded by legal-domain evidence and cannot hijack the service-readiness milestone.

## Must-Haves

- At most two candidate models are considered; broad model search is explicitly out of scope.
- Candidate evaluation uses deployment-scoped endpoints/runtime configuration, not the `/v1/embeddings` request `model` field as a selector.
- The legal corpus artifact excludes raw legal text and secret material, records corpus hash, selected document/query counts, model IDs, dimensions, cache namespaces, metrics or stop reasons, and one final outcome: keep_current, reject_candidate, or defer_candidate.
- Cross-model cosine/top-1 parity thresholds are not used as replacement criteria for different models; retrieval metrics and availability/operational compatibility drive the quick gate.
- Verification includes executable artifact validation and leak/bounds checks.

## Proof Level

- This slice proves: operational evidence when candidate endpoints are available; otherwise bounded contract/truthful-deferral evidence with executable artifact validation. Real runtime is required for any candidate that reaches the legal evidence gate; no human/UAT is required.

## Integration Closure

Consumes S01 same-host contract, current legal corpus/evaluator patterns, M039 legal baseline evidence, and Docker/runtime configuration. Introduces only evaluation tooling and a sanitized S03 evidence artifact; it does not add per-request model routing or change fd production runtime. S04 remains responsible for final TEI-vs-ONNX recommendation synthesis.

## Verification

- S03 must make candidate status inspectable through the artifact: endpoint URL labels, runtime model/dimensions/cache namespace from `/health` when available, smoke embedding dimension results, stop reason on unavailable/failed candidates, and redaction status. Failure visibility is file-based and verifier-based, not runtime API changes.

## Tasks

- [x] **T01: Add bounded legal model gate tooling** `est:3h`
  Expected executor skills for task-plan frontmatter: grill-me, sentence-transformers, api-design, write-docs, verify-before-complete.
  - Files: `tools/evaluate_legal_model_quick_gate.py`, `tools/verify_legal_model_quick_gate_artifact.py`
  - Verify: python3 -m py_compile tools/evaluate_legal_model_quick_gate.py tools/verify_legal_model_quick_gate_artifact.py
python3 tools/verify_legal_model_quick_gate_artifact.py --self-test
python3 tools/evaluate_legal_model_quick_gate.py --dry-run --corpus tests/44-FZ-2026-articles.jsonl --output benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md --baseline-model deepvk/USER-bge-m3 --baseline-runtime-label tei-default --baseline-dimensions 1024 --candidate-model BAAI/bge-m3 --candidate-runtime-label candidate-bge-m3 --candidate-dimensions 1024 --candidate-cache-namespace m040-s03-candidate-bge-m3 --max-docs 32 --max-title-queries 8 --max-self-queries 8
python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md --max-candidates 2

- [x] **T02: Run bounded candidate gate and publish S03 artifact** `est:3h`
  Expected executor skills for task-plan frontmatter: sentence-transformers, api-design, write-docs, verify-before-complete.
  - Files: `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`
  - Verify: python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03.md --max-candidates 2
python3 -m py_compile tools/evaluate_legal_model_quick_gate.py tools/verify_legal_model_quick_gate_artifact.py

## Files Likely Touched

- tools/evaluate_legal_model_quick_gate.py
- tools/verify_legal_model_quick_gate_artifact.py
- benchmark-results/fd-legal-model-quick-gate-m040-s03.md
