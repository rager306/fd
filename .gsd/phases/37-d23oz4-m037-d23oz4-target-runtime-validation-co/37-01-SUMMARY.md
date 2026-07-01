---
id: S01
parent: M037-d23oz4
milestone: M037-d23oz4
provides:
  - Actionable target-runtime acceptance policy before ONNX promotion.
requires:
  []
affects:
  - S02
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt
key_decisions: []
patterns_established:
  - Python drivers can collect evidence only when they target actual runtime endpoints; Python-only verification is not runtime acceptance.
observability_surfaces:
  - Outcome artifact records the boundary and required target-runtime gates.
drill_down_paths:
  - .gsd/milestones/M037-d23oz4/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M037-d23oz4/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M037-d23oz4/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M037-d23oz4/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T10:17:55.214Z
blocker_discovered: false
---

# S01: Target runtime validation contract

**Defined the target-runtime validation contract for Go ONNX and future Rust paths.**

## What Happened

S01 documented the target-runtime validation boundary requested by the user. The contract makes Python helper checks setup/provenance evidence only, requires Go fd API/package gates for any new/regenerated ONNX artifact, and requires equivalent independent gates for any future Rust backend.

## Verification

Manifest JSON, provisioning dry-run, artifact verifier, export contract verifier, actionlint, custom contract/leak checks, and GitNexus detect all passed.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

- Production acceptance of ONNX artifacts requires target runtime validation through the intended runtime, not only Python helper checks.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

This is a validation contract only; it does not execute new Go/Rust gates.

## Follow-ups

S02 should record the decision, run full guardrails, validate/complete milestone, checkpoint, commit, and reindex.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — Adds source_contract.target_runtime_validation contract.
- `docs/onnx-artifacts/PROVISIONING.md` — Documents Python-helper boundary, Go runtime gates, future Rust rule, Redis isolation, and blockers.
- `docs/onnx-artifacts/README.md` — Summarizes target-runtime validation boundary.
- `benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt` — Outcome artifact for M037 S01.
