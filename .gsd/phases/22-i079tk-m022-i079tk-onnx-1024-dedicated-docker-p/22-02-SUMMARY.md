---
id: S02
parent: M022-i079tk
milestone: M022-i079tk
provides:
  - A truthful hosted-CI boundary for ONNX packaging work.
requires:
  []
affects:
  - Default Go Quality workflow now covers ONNX contract metadata and binary hygiene without provisioning artifacts.
key_files:
  - .github/workflows/go-quality.yml
  - docs/onnx-artifacts/README.md
  - .gsd/DECISIONS.md
key_decisions:
  - D020: CI runs artifact-free ONNX contract checks now; full ONNX image CI deferred until external artifact provisioning/cache exists.
patterns_established:
  - Use `--allow-missing` only for CI metadata checks, never runtime readiness evidence.
  - Keep full ONNX image CI separate from default Go Quality until artifacts are provisioned externally.
observability_surfaces:
  - GitHub Actions steps for artifact contract metadata and binary hygiene; README CI boundary; D020 decision.
drill_down_paths:
  - .gsd/milestones/M022-i079tk/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M022-i079tk/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M022-i079tk/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T10:48:30.839Z
blocker_discovered: false
---

# S02: CI artifact provisioning boundary

**S02 added CI-safe ONNX contract checks and documented the full image CI provisioning blocker.**

## What Happened

S02 added honest CI coverage for the ONNX packaging boundary. The workflow now triggers for ONNX packaging docs/tooling changes, validates manifest metadata with `--allow-missing`, and fails if ONNX/native/runtime binaries are tracked. Documentation and D020 state that full ONNX image CI requires external artifact provisioning/cache before it can be truthfully enabled. Closure verification passed across actionlint, CI-equivalent checks, Go tests, lint, tagged tests, Docker builds, ONNX image smoke, binary hygiene, cleanup, and GitNexus scope.

## Verification

All S02 and milestone closure verification passed.

## Requirements Advanced

- onnx-ci-contract-check — Added CI checks for ONNX artifact metadata and binary hygiene without requiring local artifacts.

## Requirements Validated

- m022-ci-boundary — actionlint and local CI-equivalent verifier/binary-hygiene commands passed.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Hosted CI still does not build or run the ONNX image. That is intentional until artifact provisioning is implemented.

## Follow-ups

Future milestone can add full ONNX hosted CI only after selecting/provisioning an artifact store/cache for the ONNX model, native tokenizer static lib, tokenizer JSON, and ONNX Runtime shared library.

## Files Created/Modified

- `.github/workflows/go-quality.yml` — Added CI-safe ONNX contract and binary hygiene checks.
- `docs/onnx-artifacts/README.md` — Documented CI boundary and full image provisioning requirements.
- `.gsd/DECISIONS.md` — Decision register updated with D020.
