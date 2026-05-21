---
id: S01
parent: M030-3kdha1
milestone: M030-3kdha1
provides:
  - Code-level remediation for M028 LOW-3 and LOW-4.
requires:
  []
affects:
  - S02 documentation and closure
  - Future immutable artifact source selection
key_files:
  - api/embed/onnx_manifest.go
  - tools/provision_onnx_artifacts.py
  - tools/verify_onnx_artifacts.py
  - tools/build_onnx_image.sh
key_decisions:
  - Approved roots are `.gsd/runtime/onnx`, `.gsd/runtime/tokenizers`, `.gsd/runtime/onnxruntime`, and `tei-models`.
  - Default diagnostics prefer repo-relative/basename-safe path display instead of absolute host paths.
patterns_established:
  - Security path policies should align test fixtures to real approved artifact roots rather than using arbitrary temp artifact paths.
  - Prefer repo-relative/basename-safe diagnostics by default; keep deeper local context outside default logs.
observability_surfaces:
  - Safer artifact path errors and tool JSON output; outcome pending S02.
drill_down_paths:
  - .gsd/milestones/M030-3kdha1/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M030-3kdha1/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M030-3kdha1/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T05:31:20.975Z
blocker_discovered: false
---

# S01: Artifact path policy and safe diagnostics

**S01 remediated path-root and default path-disclosure risks in ONNX artifact tooling.**

## What Happened

S01 implemented path-root policy and safe diagnostics across Go manifest validation and Python artifact tooling. Repo-external, absolute, traversal, and unapproved-root manifest paths are rejected. Existing approved runtime/model roots remain supported. Default dry-run/verifier output no longer emits absolute repo root, and build-script missing diagnostics name env keys rather than configured absolute values. Guardrails passed after aligning test fixtures to the approved layout.

## Verification

All S01 verification passed.

## Requirements Advanced

- onnx-path-security — Implemented path-root policy and safer default diagnostics for ONNX artifacts.

## Requirements Validated

- m030-s01-guardrails — All S01 guardrails passed after fixture correction.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Initial full Go test run exposed test fixture mismatch with the new root policy; fixture was updated and rerun passed.

## Known Limitations

Existing manifest path policy still assumes tracked manifests are trusted; immutable artifact source selection remains future work.

## Follow-ups

S02 should document that M028 LOW-3 and LOW-4 are remediated for default diagnostics/path roots, while immutable source selection and hosted workflow proof remain open.

## Files Created/Modified

- `api/embed/onnx_manifest.go` — Go manifest path policy and safe display.
- `api/embed/onnx_manifest_test.go` — Go manifest path policy tests.
- `api/main_test.go` — Config test fixture aligned to approved artifact root.
- `tools/provision_onnx_artifacts.py` — Provisioning path policy and safe display.
- `tools/verify_onnx_artifacts.py` — Verifier path policy and safe display.
- `tools/build_onnx_image.sh` — Build script missing-artifact diagnostics avoid printing configured absolute paths.
