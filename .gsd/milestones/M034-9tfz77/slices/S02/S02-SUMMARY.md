---
id: S02
parent: M034-9tfz77
milestone: M034-9tfz77
provides:
  - Safe workflow input contract for future hosted ONNX packaging proof.
requires:
  []
affects:
  []
key_files:
  - .github/workflows/onnx-packaging.yml
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-workflow-input-alignment-m034-s02.txt
key_decisions:
  - D032: `onnx_runtime_sha256` is optional workflow override; manifest `source_contract.onnx_runtime.library_sha256` is the default checksum source.
patterns_established:
  - Manual workflow inputs should use tracked source-contract checksums by default and treat workflow dispatch as an explicit external action.
observability_surfaces:
  - Workflow validation log message for manifest-derived runtime sha.
  - Outcome artifact with input contract and verification evidence.
drill_down_paths:
  - .gsd/milestones/M034-9tfz77/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M034-9tfz77/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M034-9tfz77/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T07:59:58.933Z
blocker_discovered: false
---

# S02: Hosted workflow input contract documentation

**Documented and verified the hosted ONNX workflow input contract.**

## What Happened

S02 documented the manual hosted workflow input contract, recorded the outcome and D032, and ran final verification. The workflow/docs now clearly state required and optional inputs, the manifest-derived runtime checksum behavior, the prohibition on signed/plain secret URLs, and the need for explicit approval before any workflow dispatch.

## Verification

S02 verification passed: actionlint, provisioning/export checks, Go tests, lint, tagged tests, Docker build, docs leak checks, tracked binary hygiene, GitNexus detect, background process check, and port check all passed.

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

No hosted workflow was dispatched. Exact ONNX model binary immutable source is still missing. ONNX remains opt-in experimental.

## Follow-ups

Next gate remains exact ONNX model binary source. Do not push or dispatch the workflow until the user explicitly approves and safe immutable source inputs exist.

## Files Created/Modified

- `.github/workflows/onnx-packaging.yml` — Manual workflow runtime sha input is optional override and provisioning uses manifest sha when omitted.
- `docs/onnx-artifacts/PROVISIONING.md` — Documents manual hosted workflow input contract.
- `docs/onnx-artifacts/README.md` — Links README CI boundary to the input contract.
- `benchmark-results/fd-onnx-workflow-input-alignment-m034-s02.txt` — Outcome artifact for M034.
- `.gsd/DECISIONS.md` — Decision D032 added by GSD.
