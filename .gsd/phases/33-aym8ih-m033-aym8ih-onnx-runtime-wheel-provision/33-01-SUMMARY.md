---
id: S01
parent: M033-aym8ih
milestone: M033-aym8ih
provides:
  - Provisioning support needed before future hosted ONNX Runtime wheel proof.
requires:
  []
affects:
  - S02
key_files:
  - tools/provision_onnx_artifacts.py
key_decisions:
  - Provisioning infers zip/wheel extraction for `.zip`/`.whl` runtime sources using `source_contract.onnx_runtime.library_member`; non-zip runtime sources remain direct-file copies.
patterns_established:
  - Use source extension plus manifest member metadata to choose wheel extraction while preserving direct-file fallback.
observability_surfaces:
  - Provisioning dry-run now includes manifest-derived ONNX Runtime expected sha when source_contract metadata is present.
  - Provisioning result includes `onnx_runtime` verification records when runtime source is supplied.
drill_down_paths:
  - .gsd/milestones/M033-aym8ih/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M033-aym8ih/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M033-aym8ih/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T07:16:07.521Z
blocker_discovered: false
---

# S01: ONNX Runtime wheel provisioning support

**Made the M031 ONNX Runtime wheel source candidate actionable in provisioning tooling.**

## What Happened

S01 implemented and verified safe ONNX Runtime wheel extraction in the provisioning helper. The helper now reads runtime member/size/sha from the ONNX manifest source contract, extracts the configured member from `.whl`/`.zip` sources without `extractall`, enforces member size bounds, and verifies the destination checksum. Synthetic probes proved positive extraction, missing member, oversized member, checksum mismatch, and direct-file fallback behavior.

## Verification

S01 verification passed: py_compile, synthetic positive wheel extraction, negative probes, direct-file fallback, provisioning dry-run, local artifact verifier, and export contract verifier all passed.

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

Synthetic probes used small fake wheels. The real PyPI wheel source remains a future hosted/local provisioning input; no large runtime artifact was committed.

## Follow-ups

S02 should document wheel extraction behavior and record outcome/decision. Full hosted proof still requires exact ONNX model source and explicit push/workflow approval.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py` — Adds ONNX Runtime zip/wheel member extraction support and manifest-derived runtime sha/member metadata handling.
