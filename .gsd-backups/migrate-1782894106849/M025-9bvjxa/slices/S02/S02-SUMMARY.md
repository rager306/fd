---
id: S02
parent: M025-9bvjxa
milestone: M025-9bvjxa
provides:
  - Operational contract for hosted CI and rollout planning.
requires:
  []
affects:
  - S03 hosted ONNX CI skeleton
key_files:
  - docs/onnx-artifacts/OPERATIONS.md
  - .gsd/DECISIONS.md
key_decisions:
  - D023: ONNX rollout remains staged and opt-in; TEI remains default rollback path until diagnostics/provisioning/CI/security/rollout gates pass.
patterns_established:
  - Document operational readiness separately from quality/performance evidence.
  - Rollback must be TEI/default and immediate.
  - Failure diagnostics must name artifact_id/path/expected/actual values without secrets or raw input text.
observability_surfaces:
  - Operations runbook, required diagnostic fields, health metadata recommendation, D023 decision.
drill_down_paths:
  - .gsd/milestones/M025-9bvjxa/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M025-9bvjxa/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M025-9bvjxa/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T11:44:13.412Z
blocker_discovered: false
---

# S02: Operational diagnostics and rollout contract

**S02 established the ONNX operational diagnostics and rollout contract.**

## What Happened

S02 documented the ONNX operational diagnostics and rollout/rollback contract. It defines startup preflight, actionable failure messages, health and observability expectations, rollout stages, and rollback steps. D023 records the policy that ONNX remains opt-in and TEI remains the rollback/default path until operational gates are implemented and verified. Verification confirmed required sections and binary hygiene.

## Verification

Operations docs and binary hygiene verification passed.

## Requirements Advanced

- onnx-operational-contract — Defined operational diagnostics, rollout, and rollback requirements for ONNX opt-in runtime.

## Requirements Validated

- m025-s02-operations-doc — Operations sections and README link checks passed.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Operational health detail surfaces and startup preflight endpoint/status are documented but not implemented in code yet.

## Follow-ups

S03 should add a safe manual hosted ONNX CI skeleton that uses provisioning inputs and does not become required until external artifacts exist.

## Files Created/Modified

- `docs/onnx-artifacts/OPERATIONS.md` — Operational diagnostics and rollout/rollback contract.
- `docs/onnx-artifacts/README.md` — README links operations contract.
- `.gsd/DECISIONS.md` — Decision register updated with D023.
