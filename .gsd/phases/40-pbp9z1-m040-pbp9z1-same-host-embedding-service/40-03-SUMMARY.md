---
id: S03
parent: M040-pbp9z1
milestone: M040-pbp9z1
provides:
  - A bounded legal-domain candidate model quick-gate artifact for S04.
  - Executable quick-gate and artifact-validation tooling.
  - A fail-closed conclusion that alternative model replacement is deferred until runtime metadata and live legal retrieval evidence are available.
requires:
  - slice: S01
    provides: Same-host service contract requiring runtime metadata on `/health` and no silent runtime/model fallback.
affects:
  - S04
key_files:
  - tools/evaluate_legal_model_quick_gate.py
  - tools/verify_legal_model_quick_gate_artifact.py
  - benchmark-results/fd-legal-model-quick-gate-m040-s03.md
  - benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md
key_decisions:
  - Alternative model exploration is capped to two candidates: `BAAI/bge-m3` and `intfloat/multilingual-e5-large`.
  - Cross-model cosine/top-1 parity is not a valid replacement criterion for different embedding models; retrieval metrics and operational compatibility drive the gate.
  - A baseline `/health` response without runtime metadata makes live candidate comparison unavailable, so the gate must defer fail-closed.
patterns_established:
  - Use a separate artifact verifier for sanitized benchmark evidence rather than trusting evaluator output alone.
  - Represent unavailable runtime evidence with explicit stop reasons and `defer_candidate` instead of implicit success.
  - Use deployment-scoped endpoint/runtime configuration and cache namespaces for model comparisons, not the request `model` field as a selector.
observability_surfaces:
  - `benchmark-results/fd-legal-model-quick-gate-m040-s03.md` records endpoint labels, runtime/health stop reasons, cache namespaces, model dimensions, corpus hash/counts, redaction status, and final verdict.
  - `tools/verify_legal_model_quick_gate_artifact.py` provides executable validation for candidate bounds, redaction, metadata consistency, and verdict shape.
drill_down_paths:
  - .gsd/milestones/M040-pbp9z1/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M040-pbp9z1/slices/S03/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-22T08:16:10.141Z
blocker_discovered: false
---

# S03: Bounded legal model quick gate

**S03 bounded alternative embedding-model exploration to two legal-domain candidates and published a sanitized fail-closed quick-gate artifact for S04.**

## What Happened

S03 added executable tooling for a legal-domain model quick gate and used it to publish the canonical artifact at `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`. The evaluator and verifier enforce the slice boundary: at most two candidates, deployment-scoped endpoints/configuration instead of request-level model switching, sanitized corpus metadata only, model/runtime/dimension/cache namespace disclosure, and retrieval/operational compatibility as the only valid replacement basis. The canonical run considered exactly `BAAI/bge-m3` and `intfloat/multilingual-e5-large` against the current `deepvk/USER-bge-m3` baseline. It did not run live retrieval metrics because the same-host baseline `/health` response lacked the contract-required `runtime` object, so candidate evaluation stopped fail-closed with `defer_candidate` rather than treating an uninspectable runtime as comparable. This preserves the S01 service-contract boundary and gives S04 a clear alternative-model conclusion: keep the current model for the runtime recommendation unless a future deployment restores runtime metadata and reruns the bounded legal evidence gate.

## Verification

Closeout verification was run through `gsd_exec` in run `4cd519ac-4e7d-4f99-bc52-6e8f2af50090` and passed. It compiled `tools/evaluate_legal_model_quick_gate.py` and `tools/verify_legal_model_quick_gate_artifact.py`, ran the verifier self-test, regenerated and validated the dry-run artifact, validated the canonical S03 artifact with `--max-candidates 2`, and performed schema/content checks confirming candidate_count=2, baseline `deepvk/USER-bge-m3`, candidates `BAAI/bge-m3` and `intfloat/multilingual-e5-large`, unique isolated `m040-s03-candidate-*` cache namespaces, sanitized redaction status, fail-closed operational compatibility, and final verdict `defer_candidate`.

## Requirements Advanced

- R001 — S03 preserved the legal-domain quality gate by refusing to compare or recommend candidate models without runtime metadata and legal evidence.
- R006 — S03 supplies the alternative-model portion of the final TEI-vs-ONNX evidence envelope for S04.
- R008 — S03 implemented and validated a bounded two-candidate legal-domain quick gate.

## Requirements Validated

- R008 — Canonical artifact `benchmark-results/fd-legal-model-quick-gate-m040-s03.md` plus verifier checks prove candidate scope is capped at two, legal corpus evidence is sanitized, and candidate replacement is deferred fail-closed without live legal retrieval evidence.

## New Requirements Surfaced

- None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Live legal retrieval metrics were not run because the current same-host baseline `/health` response lacks the contract-required `runtime` object. The slice correctly rendered a fail-closed `defer_candidate` artifact instead of continuing with uninspectable runtime metadata.

## Known Limitations

Candidate quality remains deferred, not accepted or rejected by retrieval metrics. A future live comparison requires restoring/deploying `/health.runtime` metadata for the baseline and candidate endpoints, then rerunning the bounded gate.

## Follow-ups

S04 should treat alternative-model replacement as deferred and base the runtime recommendation on the current `deepvk/USER-bge-m3` service evidence unless runtime metadata is restored and the S03 gate is rerun.

## Files Created/Modified

- `tools/evaluate_legal_model_quick_gate.py` — Added bounded legal model quick-gate evaluator.
- `tools/verify_legal_model_quick_gate_artifact.py` — Added executable artifact validator and self-tests for bounds, redaction, verdict, and metadata consistency.
- `benchmark-results/fd-legal-model-quick-gate-m040-s03.md` — Published canonical sanitized S03 evidence artifact with two candidates and `defer_candidate` verdict.
- `benchmark-results/fd-legal-model-quick-gate-m040-s03-dry-run.md` — Rendered dry-run validation artifact used by tooling verification.
- `.gsd/REQUIREMENTS.md` — Regenerated by requirement update for R008 validation.
