---
id: S02
parent: M020-mvkq4d
milestone: M020-mvkq4d
provides:
  - Validated artifact contract and next packaging gate recommendation.
requires:
  []
affects:
  - Future Docker/CI packaging milestone
key_files:
  - .gsd/DECISIONS.md
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
key_decisions:
  - D018: existing dynamic-axis manifest is the experimental ONNX 1024 runtime contract; production remains blocked on Docker/CI packaging and artifact provisioning.
patterns_established:
  - Do not create a new artifact manifest unless there is a new binary/checksum.
  - An experimental runtime contract can advance packaging work without authorizing production default changes.
observability_surfaces:
  - Decision register entry D018 and manifest validation evidence.
drill_down_paths:
  - .gsd/milestones/M020-mvkq4d/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M020-mvkq4d/slices/S02/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T10:03:14.287Z
blocker_discovered: false
---

# S02: Contract validation and next gate

**S02 recorded and verified the ONNX 1024 artifact contract decision.**

## What Happened

S02 recorded D018 and ran fresh closure verification. The ONNX 1024 contract is valid JSON, includes expected fields and evidence links, preserves `production_default=false`, and no ONNX/native binaries are tracked. Tests, lint, tagged tokenizer tests, runtime cleanup, and GitNexus scope checks passed.

## Verification

Fresh verification passed: manifest contract validation, binary hygiene, Go tests, lint, tagged tests, runtime cleanup, and GitNexus scope.

## Requirements Advanced

- onnx-1024-artifact-contract — Moved ONNX 1024 from undocumented local runtime contract to tracked experimental artifact contract.

## Requirements Validated

- m020-contract-closure — D018 saved and fresh verification passed: manifest contract validation, tracked binaries 0, tests/lint/tagged checks pass.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Runtime startup code does not yet enforce the 1024 contract. Docker/CI packaging and external artifact provisioning remain unimplemented.

## Follow-ups

Plan a future milestone for Docker/CI packaging and artifact provisioning: supply ONNX binary and native tokenizer without committing binaries, verify checksums, run tagged tests in packaged environment, and rerun quality/performance gates.

## Files Created/Modified

- `.gsd/DECISIONS.md` — Decision register updated with D018.
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — ONNX manifest contract validated.
