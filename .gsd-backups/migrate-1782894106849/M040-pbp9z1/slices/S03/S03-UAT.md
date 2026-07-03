# S03: Bounded legal model quick gate — UAT

**Milestone:** M040-pbp9z1
**Written:** 2026-05-22T08:16:10.141Z

# S03 UAT: Bounded legal model quick gate

## UAT Type
Automated artifact and contract verification; no human interactive service test is required for this slice.

## Preconditions
1. Work from `/root/fd`.
2. `tools/evaluate_legal_model_quick_gate.py`, `tools/verify_legal_model_quick_gate_artifact.py`, `tests/44-FZ-2026-articles.jsonl`, and `benchmark-results/fd-legal-model-quick-gate-m040-s03.md` exist.
3. Candidate endpoints may be unavailable; if the baseline or candidates are uninspectable, the expected behavior is a truthful `defer_candidate` stop reason, not a silent pass.

## Steps
1. Compile the quick-gate tools:
   `python3 -m py_compile tools/evaluate_legal_model_quick_gate.py tools/verify_legal_model_quick_gate_artifact.py`
2. Run verifier self-tests:
   `python3 tools/verify_legal_model_quick_gate_artifact.py --self-test`
3. Generate a bounded dry-run artifact with one candidate:
   `python3 tools/evaluate_legal_model_quick_gate.py --dry-run --corpus tests/44-FZ-2026-articles.jsonl --output benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md --baseline-model deepvk/USER-bge-m3 --baseline-runtime-label tei-default --baseline-dimensions 1024 --candidate-model BAAI/bge-m3 --candidate-runtime-label candidate-bge-m3 --candidate-dimensions 1024 --candidate-cache-namespace m040-s03-candidate-bge-m3 --max-docs 32 --max-title-queries 8 --max-self-queries 8`
4. Validate the dry-run artifact:
   `python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md --max-candidates 2`
5. Validate the canonical S03 artifact:
   `python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03.md --max-candidates 2`
6. Inspect the canonical artifact and confirm it records: corpus SHA-256 and counts, no raw legal text, two or fewer candidates, deployment endpoint labels, model IDs, dimensions, cache namespaces, runtime/health stop reasons, and a final verdict.

## Expected Outcomes
- All commands exit 0.
- The canonical artifact validates with `--max-candidates 2`.
- Candidate scope is capped to `BAAI/bge-m3` and `intfloat/multilingual-e5-large`.
- The artifact uses sanitized hashes/counts and explicitly states raw legal corpus text is excluded.
- The artifact does not use cross-model cosine/top-1 parity as a replacement criterion.
- The current canonical outcome is `defer_candidate` because the baseline `/health` lacks the required runtime metadata.

## Edge Cases
- If more than two candidates are configured, evaluation/verification must fail closed.
- If `/health` lacks runtime metadata, the artifact must contain a stop reason and defer candidate comparison.
- If raw legal text or secret-like material appears in the artifact, verification must fail.
- If candidate cache namespaces collide or are omitted, closeout should fail until fixed.

## Not Proven By This UAT
- It does not prove a candidate model outperforms or matches `deepvk/USER-bge-m3` on live legal retrieval metrics.
- It does not prove candidate service deployment readiness.
- It does not change the production `/v1/embeddings` API contract or add per-request model routing.
