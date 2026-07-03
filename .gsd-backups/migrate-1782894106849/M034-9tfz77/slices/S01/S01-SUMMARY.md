---
id: S01
parent: M034-9tfz77
milestone: M034-9tfz77
provides:
  - Manual workflow surface that can use the M031/M033 ONNX Runtime wheel candidate without redundant sha input.
requires:
  []
affects:
  - S02
key_files:
  - .github/workflows/onnx-packaging.yml
key_decisions:
  - Workflow runtime sha input is optional when manifest source_contract provides library_sha256; CLI override remains available.
patterns_established:
  - Workflow inputs should avoid duplicating manifest checksum metadata unless an override is explicitly needed.
observability_surfaces:
  - Workflow validation message states when manifest-derived runtime sha will be used.
drill_down_paths:
  - .gsd/milestones/M034-9tfz77/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M034-9tfz77/slices/S01/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T07:53:18.229Z
blocker_discovered: false
---

# S01: Workflow runtime source alignment

**Aligned hosted ONNX workflow runtime-source input handling with provisioning behavior.**

## What Happened

S01 aligned the manual ONNX packaging workflow with M033 provisioning behavior. The workflow now treats `onnx_runtime_sha256` as an optional override and lets provisioning use manifest `source_contract.onnx_runtime.library_sha256` when omitted. Actionlint and local provisioning dry-run confirmed the alignment.

## Verification

S01 verification passed: workflow text checks, actionlint, py_compile, and provisioning dry-run alignment checks all passed.

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

No remote workflow was run. Exact ONNX model source URL/cache key is still missing.

## Follow-ups

S02 should document safe workflow dispatch input contract and explicitly state no workflow was dispatched and exact ONNX model source remains blocked.

## Files Created/Modified

- `.github/workflows/onnx-packaging.yml` — Manual workflow runtime sha input is optional override; provisioning uses manifest sha when omitted.
