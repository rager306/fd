---
id: S02
parent: M011-33b7wf
milestone: M011-33b7wf
provides:
  - `embed.ValidateONNXArtifactManifest` for S03.
  - `EMBEDDING_BACKEND` config seam with TEI default.
  - Test-covered ONNX manifest validation errors.
requires:
  []
affects:
  - S03
  - S04
key_files:
  - api/embed/onnx_manifest.go
  - api/main.go
  - api/main_test.go
  - api/embed/onnx_manifest_test.go
key_decisions:
  - TEI remains default when `EMBEDDING_BACKEND` is unset.
  - ONNX requires `ONNX_ARTIFACT_MANIFEST` and valid checksum before future loading.
  - Explicit ONNX must not silently fall back to TEI in benchmark/prototype mode.
patterns_established:
  - Manifest validation before loader.
  - Explicit opt-in runtime errors before silent fallback.
  - Pure validator package before ONNX runtime dependency.
observability_surfaces:
  - Sentinel validation errors for missing artifact, checksum mismatch, metadata mismatch, and production-default manifest.
  - Startup logs selected backend and validated ONNX artifact metadata before the temporary not-implemented exit.
drill_down_paths:
  - .gsd/milestones/M011-33b7wf/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M011-33b7wf/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M011-33b7wf/slices/S02/tasks/T03-SUMMARY.md
  - .gsd/milestones/M011-33b7wf/slices/S02/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T19:06:14.486Z
blocker_discovered: false
---

# S02: Opt in backend seam

**S02 added tested ONNX manifest validation and opt-in backend config while preserving TEI as default.**

## What Happened

S02 added the backend validation seam without changing production behavior. A new `api/embed` manifest validator parses the tracked ONNX manifest and validates file existence, size, SHA256, output metadata, dimensions, normalization expectation, and prototype-only status. Startup config now defaults to TEI, rejects invalid backend names, and requires/validates the ONNX manifest only when `EMBEDDING_BACKEND=onnx`. Since S03 has not implemented ONNX inference yet, the explicit ONNX branch validates then exits with a not-implemented error instead of silently serving TEI. Tests and lint pass.

## Verification

S02 final verification passed with fresh evidence in T04: Go tests 72 passed, lint 0 issues, Compose config OK, manifest checksum OK, and GitNexus low risk/no affected processes.

## Requirements Advanced

- future ONNX artifact distribution/checksum requirement — Introduced tracked artifact manifest validation and explicit opt-in backend config primitives.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

`EMBEDDING_BACKEND=onnx` validates the manifest but exits with a not-implemented error rather than serving TEI fallback. This is intentional to preserve benchmark integrity until S03 implements the actual loader.

## Known Limitations

ONNX inference is not implemented yet. The manifest validator hashes the full artifact at validation time, which may be acceptable for startup but should be benchmarked if startup latency matters.

## Follow-ups

S03 should implement the opt-in ONNX dense backend loader using the validated manifest. It should remove the temporary not-implemented exit path only when actual ONNX inference is wired and tested.

## Files Created/Modified

- `api/embed/onnx_manifest.go` — Pure Go ONNX artifact manifest parser and validator.
- `api/embed/onnx_manifest_test.go` — Manifest validator unit tests for valid, missing, checksum mismatch, metadata mismatch, production default, and invalid JSON cases.
- `api/main.go` — Backend config parsing and manifest validation startup seam.
- `api/main_test.go` — Backend config tests preserving TEI default and ONNX validation behavior.
