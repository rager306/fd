---
id: S01
parent: M035-7j2h6x
milestone: M035-7j2h6x
provides:
  - Actionable exact-binary source contract for future hosted ONNX packaging proof.
requires:
  []
affects:
  - S02
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt
key_decisions: []
patterns_established:
  - Future ONNX artifact source contracts can include planned key templates only if they are explicitly marked not uploaded and not usable as source URLs.
observability_surfaces:
  - Outcome artifact records exact binary identity, planned key, and remaining blocker state.
drill_down_paths:
  - .gsd/milestones/M035-7j2h6x/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M035-7j2h6x/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M035-7j2h6x/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T09:22:43.261Z
blocker_discovered: false
---

# S01: Exact binary hosting contract

**Defined the planned exact ONNX binary hosting contract while keeping source availability blocked.**

## What Happened

S01 converted the vague exact ONNX binary source blocker into a concrete hosting contract without creating a fake source. The manifest now records a `planned_not_uploaded` hosting contract with recommended key/filename, allowed/forbidden source forms, and pre-dispatch requirements. Docs and outcome mirror the same policy and preserve the production/default and external-action boundaries.

## Verification

Manifest JSON, provisioning dry-run, artifact verifier, export contract verifier, actionlint, custom contract/leak checks, and GitNexus detect all passed.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

The contract is planned only; no binary source exists yet and hosted workflow proof remains blocked.

## Follow-ups

S02 should record the decision, run final guardrails, complete the milestone, checkpoint DB, commit locally, and reindex GitNexus.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — Adds planned_not_uploaded hosting_contract under source_contract.onnx_model_artifact while keeping source_status blocked and no source_url.
- `docs/onnx-artifacts/PROVISIONING.md` — Documents exact binary size/sha, object key template, allowed/forbidden source forms, and pre-dispatch checklist.
- `docs/onnx-artifacts/README.md` — Summarizes exact-binary contract in artifact README.
- `benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt` — Outcome artifact for the exact binary hosting contract.
