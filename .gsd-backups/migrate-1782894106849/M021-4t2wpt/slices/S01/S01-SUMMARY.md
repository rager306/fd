---
id: S01
parent: M021-4t2wpt
milestone: M021-4t2wpt
provides:
  - Artifact verification script and README for Docker/CI packaging milestone.
requires:
  []
affects:
  - S02
  - Future Docker/CI packaging implementation
key_files:
  - tools/verify_onnx_artifacts.py
  - docs/onnx-artifacts/README.md
key_decisions:
  - Do not alter default Dockerfile or CI in S01; keep default TEI path independent of ONNX/native artifacts.
  - Use `tools/verify_onnx_artifacts.py` as the local fail-fast contract for artifact provisioning and checksum validation.
  - Use `--allow-missing` only for default CI contract-shape checks, never as ONNX runtime readiness evidence.
patterns_established:
  - Artifact readiness requires checksum verification against manifests, not file presence alone.
  - Default CI may check contract shape with missing artifacts, but ONNX readiness requires present verified artifacts.
observability_surfaces:
  - Verifier JSON output with artifact IDs, paths, sizes, checksums, and verified flags.
drill_down_paths:
  - .gsd/milestones/M021-4t2wpt/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M021-4t2wpt/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M021-4t2wpt/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T10:16:09.125Z
blocker_discovered: false
---

# S01: Artifact provisioning contract

**S01 added and validated the ONNX/native artifact provisioning contract without changing default Docker/CI.**

## What Happened

S01 created and validated the artifact provisioning contract. The new verifier checks tracked manifests against local ignored ONNX and native tokenizer files, verifying presence, size, checksum, production_default false, artifact.git_tracked false, and git-tracked path hygiene. The README documents required artifacts, verification commands, the validated ONNX 1024 runtime env, existing evidence, and the future Docker/CI packaging gate. No binary artifacts were tracked.

## Verification

S01 verification passed: script compile, strict artifact verification, allow-missing check, manifest JSON validation, tracked binary hygiene, and GitNexus scope.

## Requirements Advanced

- onnx-1024-packaging-contract — Defined local verification contract for ONNX/native artifact provisioning without tracking binaries.

## Requirements Validated

- m021-s01-artifact-contract — `tools/verify_onnx_artifacts.py` verifies local ONNX and native tokenizer artifacts against manifests; tracked binary check reports 0.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

No artifact download/provisioning implementation yet. No Docker/CI workflow changes yet. The verifier assumes artifacts are already present locally unless `--allow-missing` is explicitly used.

## Follow-ups

S02 should validate default Docker/build boundary and record that the next implementation gate is CI/Docker artifact provisioning, not ONNX promotion.

## Files Created/Modified

- `tools/verify_onnx_artifacts.py` — Local verifier for ONNX and native tokenizer artifacts against tracked manifests.
- `docs/onnx-artifacts/README.md` — Artifact provisioning contract and future Docker/CI gate notes.
