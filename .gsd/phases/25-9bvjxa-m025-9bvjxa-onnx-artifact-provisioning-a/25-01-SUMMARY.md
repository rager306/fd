---
id: S01
parent: M025-9bvjxa
milestone: M025-9bvjxa
provides:
  - Provisioning contract and helper for hosted CI/deploy design.
requires:
  []
affects:
  - S02 operational diagnostics and S03 hosted ONNX CI skeleton
key_files:
  - docs/onnx-artifacts/PROVISIONING.md
  - tools/provision_onnx_artifacts.py
key_decisions:
  - Missing required external sources are blockers, not defaults.
  - Provisioning helper requires explicit sources in non-dry-run mode.
  - Checksum verification is mandatory before any ONNX build/runtime use.
patterns_established:
  - Treat missing external artifact URLs as blockers.
  - Separate metadata-shape CI from runtime-readiness CI.
  - Require checksum verification before ONNX build/runtime evidence.
observability_surfaces:
  - Provisioning dry-run JSON, clear missing-source failure messages, docs failure/diagnostics contract.
drill_down_paths:
  - .gsd/milestones/M025-9bvjxa/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M025-9bvjxa/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M025-9bvjxa/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T11:40:46.081Z
blocker_discovered: false
---

# S01: Artifact provisioning contract

**S01 established the ONNX artifact provisioning contract and helper without committing binaries.**

## What Happened

S01 documented and implemented the ONNX artifact provisioning contract. It identifies all required artifacts, destination paths, checksum policy, source selection priorities, cache layout, and current hosted-CI blockers. `tools/provision_onnx_artifacts.py` now provides a dry-run plan and checksum-verified explicit-source provisioning path. Verification passed for dry-run, missing-source failure, strict local verifier, docs discoverability, and binary hygiene.

## Verification

All S01 verification passed.

## Requirements Advanced

- onnx-artifact-provisioning-contract — Defined and implemented explicit-source artifact provisioning for ONNX/native assets.

## Requirements Validated

- m025-s01-provisioning-helper — Dry-run and missing-source failure behavior passed; strict local verifier passed.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

The ONNX model external source and pinned ONNX Runtime source are still undefined. Native tokenizer source uses upstream `latest` and should be pinned or mirrored for production.

## Follow-ups

S02 should define operational diagnostics and rollout/rollback behavior; S03 should wire provisioning into safe manual hosted CI skeleton.

## Files Created/Modified

- `docs/onnx-artifacts/PROVISIONING.md` — Provisioning/cache contract for ONNX artifacts.
- `docs/onnx-artifacts/README.md` — README now links provisioning contract.
- `tools/provision_onnx_artifacts.py` — Dry-run and checksum-verified provisioning helper.
