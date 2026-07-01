---
id: S01
parent: M020-mvkq4d
milestone: M020-mvkq4d
provides:
  - Tracked ONNX 1024 runtime contract for packaging/CI milestone.
requires:
  []
affects:
  - S02
  - Future Docker/CI packaging milestone
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
key_decisions:
  - Existing manifest updated instead of creating a new artifact file because the ONNX binary is unchanged and has dynamic sequence axes.
  - Export sequence length 128 is preserved as provenance; validated runtime max sequence length 1024 is documented separately.
  - Manifest status is experimental/local validated and production_default remains false.
patterns_established:
  - Separate export provenance from validated runtime contract.
  - Quality/performance metadata must not imply production readiness while packaging gates remain open.
observability_surfaces:
  - Manifest fields for validated runtime env, evidence artifacts, failure contract, and future gates.
drill_down_paths:
  - .gsd/milestones/M020-mvkq4d/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M020-mvkq4d/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M020-mvkq4d/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T10:01:05.423Z
blocker_discovered: false
---

# S01: ONNX 1024 runtime contract

**S01 made the ONNX 1024 runtime contract explicit in tracked metadata.**

## What Happened

S01 updated and validated the ONNX manifest contract. The manifest now makes the 1024 runtime contract explicit, links M018 quality and M019 performance artifacts, preserves export provenance at sequence length 128, and states remaining gates before production. Binary hygiene checks confirm no ONNX or native tokenizer binaries are tracked.

## Verification

S01 verification passed: JSON/field checks, evidence links, and tracked binary hygiene checks all passed.

## Requirements Advanced

- onnx-1024-artifact-contract — Documented the ONNX 1024 runtime contract and evidence artifacts while preserving experimental status.

## Requirements Validated

- m020-s01-manifest-contract — `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` has `runtime.validated_max_sequence_length=1024`, M018/M019 evidence links, `production_default=false`, and no tracked binaries.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

The manifest contract is documentation/metadata only. Runtime startup code does not yet enforce `validated_max_sequence_length`, and Docker/CI packaging is not implemented.

## Follow-ups

S02 should record the contract decision and validate that next gate is Docker/CI packaging plus artifact provisioning. Future implementation should consume this contract in startup validation where feasible.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — Tracked ONNX manifest now records validated runtime max sequence length 1024, quality/performance evidence, and remaining production gates.
