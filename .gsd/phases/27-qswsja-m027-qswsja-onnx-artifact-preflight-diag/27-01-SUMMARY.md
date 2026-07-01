---
id: S01
parent: M027-qswsja
milestone: M027-qswsja
provides:
  - Code-level preflight diagnostics for tokenizer JSON, optional runtime library sha, and configured provider support.
requires:
  []
affects:
  - S02 documentation and closure
  - Future security/logging review
key_files:
  - api/embed/onnx_manifest.go
  - api/main.go
  - api/handlers/health.go
key_decisions:
  - Runtime library sha verification is explicit opt-in through ONNX_RUNTIME_SHA256.
  - Current configured provider support is CPUExecutionProvider only.
patterns_established:
  - Use manifest metadata for tokenizer preflight; keep runtime library hash opt-in until runtime artifact source is tracked.
  - Validate configured provider support honestly; do not claim provider enumeration when unavailable.
observability_surfaces:
  - Startup errors for tokenizer mismatch, runtime sha mismatch, unsupported provider; safe /health provider and verification flags; safe startup log flags.
drill_down_paths:
  - .gsd/milestones/M027-qswsja/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M027-qswsja/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M027-qswsja/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T12:40:28.852Z
blocker_discovered: false
---

# S01: Startup artifact and provider preflight

**S01 implemented and verified ONNX tokenizer/runtime/provider preflight diagnostics.**

## What Happened

S01 implemented tokenizer JSON checksum/size preflight using manifest metadata, optional ONNX Runtime shared library sha verification via `ONNX_RUNTIME_SHA256`, and provider validation for current CPU-only Go ONNX startup. It extended safe health metadata with provider/tokenizer/runtime verification fields and added tests for mismatch and success cases. Guardrails passed across default tests, lint, tagged tests, Docker, workflows, scripts, hygiene, and cleanup.

## Verification

All S01 guardrails passed.

## Requirements Advanced

- onnx-operational-preflight — Implemented tokenizer/runtime/provider preflight diagnostics for ONNX startup.

## Requirements Validated

- m027-s01-guardrails — All executable S01 checks passed.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Provider diagnostics validate configured provider support only; runtime provider enumeration remains future work because current startup path has no enumeration surface.

## Known Limitations

Runtime provider enumeration and deeper path redaction/security review are still future work.

## Follow-ups

S02 should update operations docs/outcome and record a decision that runtime sha is explicit opt-in and provider diagnostics are CPU-only configured-provider validation.

## Files Created/Modified

- `api/embed/onnx_manifest.go` — Manifest tokenizer.json source metadata parsing and validation exposure.
- `api/embed/onnx_manifest_test.go` — Manifest metadata tests.
- `api/main.go` — Startup preflight for tokenizer JSON, runtime sha, provider, and health population.
- `api/main_test.go` — Startup preflight tests.
- `api/handlers/health.go` — Provider/tokenizer/runtime verification fields in safe health metadata.
- `api/handlers/health_test.go` — Health metadata safety tests.
