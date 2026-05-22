---
id: T01
parent: S03
milestone: M040-pbp9z1
key_files:
  - tools/evaluate_legal_model_quick_gate.py
  - tools/verify_legal_model_quick_gate_artifact.py
  - benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md
key_decisions:
  - Bound candidate replacement acceptance to retrieval metrics and operational compatibility; cross-model cosine/parity is explicitly non-applicable.
  - Dry-run availability-only runs truthfully render defer_candidate with phase-specific stop reasons instead of implying readiness.
  - Artifact validation is separate from the evaluator and includes self-contained negative fixtures for redaction, secret-like patterns, candidate bounds, verdict, and metadata consistency.
duration: 
verification_result: passed
completed_at: 2026-05-22T08:11:53.650Z
blocker_discovered: false
---

# T01: Added bounded legal model quick-gate tooling, artifact validation, and a sanitized dry-run artifact for S03.

**Added bounded legal model quick-gate tooling, artifact validation, and a sanitized dry-run artifact for S03.**

## What Happened

Created an additive quick-gate CLI at tools/evaluate_legal_model_quick_gate.py that reuses the existing legal corpus loading, deterministic selection, hashing/redaction, embedding dimension checks, batch embedding flow, and retrieval metric patterns. The tool enforces a maximum of two candidates, supports dry-run availability-only artifacts, requires explicit baseline/candidate runtime labels, model IDs, dimensions, and cache namespaces, and uses live /health plus smoke embedding probes before legal corpus metrics when endpoints are configured. It renders final outcomes as keep_current, reject_candidate, or defer_candidate, and explicitly marks cross-model cosine/parity as not applicable for replacement-model acceptance.

Created tools/verify_legal_model_quick_gate_artifact.py to validate artifact structure, candidate count, required metadata, verdict consistency, stop-reason/metrics consistency, redaction status, absence of raw-text sections, absence of obvious secret/token patterns, and the cross-model cosine non-applicability statement. Its --self-test mode uses inline/temp fixtures and covers pass plus negative cases for too many candidates, missing metadata, raw_text_logged not false, secret-like patterns, missing verdict, and cross-model cosine acceptance claims.

Rendered benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md with the planned BAAI/bge-m3 candidate. Because this was a dry-run, baseline and candidate health/smoke/metrics are intentionally deferred with phase-specific stop reasons and no raw corpus text.

## Verification

Verified both new tools compile, the verifier self-test passes, the quick-gate dry-run renders the required sanitized artifact, the artifact verifier accepts that artifact with max two candidates, and the evaluator rejects more than two candidates with a fail-closed setup stop reason. Ran GitNexus change-scope check via local CLI because the dedicated tool was not exposed in this session; the CLI reported no detected graph changes for repo fd, consistent with this additive tooling/artifact work.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/evaluate_legal_model_quick_gate.py tools/verify_legal_model_quick_gate_artifact.py` | 0 | ✅ pass | 59ms |
| 2 | `python3 tools/verify_legal_model_quick_gate_artifact.py --self-test` | 0 | ✅ pass | 60ms |
| 3 | `python3 tools/evaluate_legal_model_quick_gate.py --dry-run --corpus tests/44-FZ-2026-articles.jsonl --output benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md --baseline-model deepvk/USER-bge-m3 --baseline-runtime-label tei-default --baseline-dimensions 1024 --candidate-model BAAI/bge-m3 --candidate-runtime-label candidate-bge-m3 --candidate-dimensions 1024 --candidate-cache-namespace m040-s03-candidate-bge-m3 --max-docs 32 --max-title-queries 8 --max-self-queries 8` | 0 | ✅ pass | 193ms |
| 4 | `python3 tools/verify_legal_model_quick_gate_artifact.py --artifact benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md --max-candidates 2` | 0 | ✅ pass | 57ms |
| 5 | `negative check: dry-run evaluator invoked with three candidate-model entries; assert exit 1 and stderr contains 'at most two candidates'` | 0 | ✅ pass | 154ms |
| 6 | `npx gitnexus detect-changes --repo fd` | 0 | ✅ pass (reported: No changes detected) | 1402ms |

## Deviations

The requested Skill tool activations could not be performed because no Skill tool is exposed in the callable tool namespace for this session. GitNexus impact analysis was not required because the change only added new files; the final GitNexus check was run through the CLI because the dedicated gitnexus_detect_changes tool was not exposed.

## Known Issues

Live endpoint behavior is implemented but not exercised in this task because the required verification path is dry-run/availability-only and no candidate service endpoint was provided. The GitNexus CLI reports no detected changes for repo fd, which may not include untracked additive files in its diff view; verification evidence therefore relies primarily on compile/self-test/dry-run/artifact validation.

## Files Created/Modified

- `tools/evaluate_legal_model_quick_gate.py`
- `tools/verify_legal_model_quick_gate_artifact.py`
- `benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md`
