---
id: S01
parent: M036-o0hewj
milestone: M036-o0hewj
provides:
  - Actionable no-upload alternative path to resolve the ONNX source blocker.
requires:
  []
affects:
  - S02
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt
key_decisions: []
patterns_established:
  - Regenerated-export proof must include both artifact regeneration evidence and downstream quality/performance/package gates, not just metadata checks.
observability_surfaces:
  - Outcome artifact records current evidence boundary and future acceptance gates.
drill_down_paths:
  - .gsd/milestones/M036-o0hewj/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M036-o0hewj/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M036-o0hewj/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M036-o0hewj/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T09:37:34.559Z
blocker_discovered: false
---

# S01: Reproducible export contract

**Defined the planned reproducible-export workflow contract without claiming proof.**

## What Happened

S01 defined the no-upload reproducible-export workflow contract. It documents pinned inputs, current verifier claim boundary, acceptance gates, success/failure interpretation, and explicit non-actions. The existing M032 verifier remains limited to existing-artifact contract verification.

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

This is planned_not_proven; no ONNX export was regenerated, no hosted workflow ran, and source blocker remains unresolved.

## Follow-ups

S02 should record decision, run full guardrails, validate/complete milestone, checkpoint, commit, and reindex.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — Adds source_contract.reproducible_export_workflow planned_not_proven contract.
- `docs/onnx-artifacts/PROVISIONING.md` — Documents reproducible-export workflow contract and acceptance gates.
- `docs/onnx-artifacts/README.md` — Summarizes no-upload reproducible-export alternative.
- `benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt` — Outcome artifact for M036 S01.
