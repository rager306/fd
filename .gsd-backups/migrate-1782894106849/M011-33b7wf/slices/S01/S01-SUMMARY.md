---
id: S01
parent: M011-33b7wf
milestone: M011-33b7wf
provides:
  - Artifact manifest for S02/S03.
  - Failure contract for missing/checksum mismatch/metadata mismatch.
  - No-production-default boundary for ONNX prototype work.
requires:
  []
affects:
  - S02
  - S03
  - S04
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - .gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md
key_decisions:
  - Use tracked JSON manifest for artifact identity.
  - Keep ONNX binary ignored and local.
  - Fail fast on explicit ONNX artifact validation errors; do not silently fall back during benchmark evidence collection.
patterns_established:
  - Manifest before loader.
  - Checksum before ONNX Runtime load.
  - Explicit ONNX validation failure before silent fallback.
observability_surfaces:
  - `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` exposes artifact_id, path, size, hash, expected metadata, and failure contract.
  - S01 research defines startup/load diagnostic expectations for missing artifact and checksum mismatch.
drill_down_paths:
  - .gsd/milestones/M011-33b7wf/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M011-33b7wf/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M011-33b7wf/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T18:59:05.594Z
blocker_discovered: false
---

# S01: Artifact manifest and checksum contract

**S01 established a tracked checksum manifest and validation contract for the local ONNX prototype artifact without committing the binary or changing runtime defaults.**

## What Happened

S01 inspected M010 artifact metadata, wrote a tracked ONNX manifest, and documented the artifact validation contract. The manifest records model/source hashes, ONNX path/size/SHA256, runtime input/output metadata, dependency pins, validation evidence, and failure expectations while leaving the 1.43GB ONNX binary ignored. The research artifact defines the validation order and confirms that ONNX remains prototype-only and non-default.

## Verification

S01 verification passed: manifest parsed and matched local artifact hash/size; git ignored status confirmed runtime artifact is not staged; S01 research includes missing artifact/checksum mismatch/no-production-default contract.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

- Future production ONNX integration needs external artifact storage or download procedure with checksum verification.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None. S01 did not touch production runtime code.

## Known Limitations

The manifest references a local `.gsd/runtime` artifact, not a production download location. Artifact distribution remains a future gate.

## Follow-ups

S02 should implement manifest/config validation and tests before any ONNX Runtime loader is introduced. Validation should cover missing artifact, checksum mismatch, metadata mismatch, and TEI default unaffected behavior.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — Tracked ONNX artifact manifest for the M010 FP32 dense candidate.
- `.gsd/milestones/M011-33b7wf/slices/S01/S01-RESEARCH.md` — Artifact validation contract and downstream implementation guidance.
