---
sliceId: S03
uatType: artifact-driven
verdict: PASS
date: 2026-05-22T08:17:00.000Z
---

# UAT Result — S03

## Checks

| Check | Mode | Result | Notes |
|-------|------|--------|-------|
| Preconditions: quick-gate tools, corpus, and canonical artifact exist and are usable from `/root/fd`. | artifact | PASS | Verified implicitly by successful compile, evaluator, verifier, and artifact inspection commands in `gsd_exec` runs `b8e08448-12a1-4ebe-be8d-5f24d0fa153c` and `f5329929-4650-4a0c-851e-1f33f2ab9cfb`. |
| Compile the quick-gate tools. | runtime | PASS | `python3 -m py_compile tools/evaluate_legal_model_quick_gate.py tools/verify_legal_model_quick_gate_artifact.py` exited 0 in run `b8e08448-12a1-4ebe-be8d-5f24d0fa153c`. |
| Run verifier self-tests. | runtime | PASS | `python3 tools/verify_legal_model_quick_gate_artifact.py --self-test` exited 0 and reported expected pass/fail cases, including too many candidates, missing metadata, raw text, secret pattern, missing verdict, and cross-model cosine acceptance rejection. |
| Generate a bounded dry-run artifact with one candidate. | runtime | PASS | `python3 tools/evaluate_legal_model_quick_gate.py --dry-run ... --candidate-model BAAI/bge-m3 ... --max-docs 32 --max-title-queries 8 --max-self-queries 8` exited 0 and wrote `benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md`; generated output showed `candidate_count: 1`, sanitized hashes/counts, no endpoint calls, and `defer_candidate`. |
| Validate the dry-run artifact. | runtime | PASS | `python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md --max-candidates 2` exited 0 with `artifact valid`. |
| Validate the canonical S03 artifact. | runtime | PASS | `python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03.md --max-candidates 2` exited 0 with `artifact valid`. |
| Inspect canonical artifact for corpus SHA-256 and counts. | artifact | PASS | Corrected structured inspection in run `f5329929-4650-4a0c-851e-1f33f2ab9cfb` parsed the JSON blocks and verified a 64-hex `corpus_sha256` plus positive `corpus_stats` counts for articles, candidate documents, clauses, parts, and title queries. |
| Inspect canonical artifact for no raw legal text. | artifact | PASS | Artifact has `raw_text_logged: false`, redaction status `sanitized_hashes_only_no_raw_legal_text`, and a statement that raw legal corpus text and smoke payload text are excluded. |
| Inspect canonical artifact for two or fewer candidates and exact candidate scope. | artifact | PASS | Structured inspection verified `candidate_count == 2` and exactly `BAAI/bge-m3` plus `intfloat/multilingual-e5-large`. |
| Inspect canonical artifact for deployment endpoint labels, model IDs, dimensions, and cache namespaces. | artifact | PASS | Structured inspection verified baseline/candidate endpoint labels, runtime labels, baseline `deepvk/USER-bge-m3`, candidate model IDs, all dimensions at 1024, and unique cache namespaces including `m040-s03-candidate-*`. |
| Inspect canonical artifact for runtime/health stop reasons and final verdict. | artifact | PASS | Structured inspection verified baseline and candidates record `/health missing runtime object` stop reasons and final verdict `defer_candidate`. |
| Confirm artifact does not use cross-model cosine/top-1 parity as a replacement criterion. | artifact | PASS | Structured inspection verified the Cross-Model Cosine/Parity block has `applicable: false` and states cross-model cosine/top-1 parity is not an acceptance metric. |

## Overall Verdict

PASS — All automatable S03 artifact-driven UAT checks passed; the canonical quick-gate artifact remains bounded to two candidates, sanitized, fail-closed, and verdicted `defer_candidate` because required runtime metadata is unavailable.

## Notes

- Evidence was gathered with `gsd_exec` per the verification-lane instruction.
- Initial run `b8e08448-12a1-4ebe-be8d-5f24d0fa153c` completed all required UAT commands successfully, then exited non-zero only because an auxiliary regex inspection was too narrow for the actual `corpus_stats` JSON shape. A corrected structured JSON inspection was run immediately afterward in `f5329929-4650-4a0c-851e-1f33f2ab9cfb` and passed.
- No human-only checks remain for this artifact-driven UAT.
