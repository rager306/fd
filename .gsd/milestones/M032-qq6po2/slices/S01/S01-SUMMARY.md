---
id: S01
parent: M032-qq6po2
milestone: M032-qq6po2
provides:
  - Executable local proof surface for S02 documentation and future hosted proof planning.
requires:
  []
affects:
  - S02
key_files:
  - tools/verify_onnx_export_contract.py
  - .gsd/milestones/M032-qq6po2/slices/S01/S01-RESEARCH.md
key_decisions:
  - The new verifier intentionally claims only existing-artifact contract verification, not regenerated export reproducibility.
patterns_established:
  - Use explicit `claim_scope` fields in verifier output to prevent rollout/reproducibility overclaims.
observability_surfaces:
  - Structured JSON verifier output with `claim_scope` and failure labels.
drill_down_paths:
  - .gsd/milestones/M032-qq6po2/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M032-qq6po2/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M032-qq6po2/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T06:57:45.005Z
blocker_discovered: false
---

# S01: Local export contract verifier

**Added a local ONNX export contract verifier with explicit proof boundaries.**

## What Happened

S01 added and verified a safe local ONNX export contract verifier. Positive verification confirms the existing local artifact matches tracked manifest, provenance, export metadata, and toolchain pins. Negative probes confirm checksum, model revision, and package drift are caught with structured sanitized errors. The proof boundary is documented to avoid overclaiming byte-for-byte reproducibility.

## Verification

S01 verification passed: py_compile and positive verifier run passed; negative tamper probes failed as expected; research boundary/leak checks passed.

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

Verifier initially expected Python inside export metadata `packages`; fixed to compare the top-level Python runtime string separately.

## Known Limitations

No full ONNX re-export was run. The ONNX model binary source blocker remains.

## Follow-ups

S02 should wire the verifier into docs/outcome and record the decision that existing-artifact contract verification is useful but not a substitute for exact binary hosting or a full reproducible export gate.

## Files Created/Modified

- `tools/verify_onnx_export_contract.py` — New local verifier for existing ONNX export contract.
- `.gsd/milestones/M032-qq6po2/slices/S01/S01-RESEARCH.md` — Verifier proof-boundary research.
