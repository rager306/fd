---
id: T01
parent: S02
milestone: M035-7j2h6x
key_files:
  - .gsd/DECISIONS.md
  - benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt
key_decisions:
  - D033: represent the exact ONNX model source blocker as planned_not_uploaded hosting contract while keeping source_status blocked and no source_url.
duration: 
verification_result: passed
completed_at: 2026-05-21T09:23:51.455Z
blocker_discovered: false
---

# T01: Recorded D033 for the planned exact ONNX binary hosting contract.

**Recorded D033 for the planned exact ONNX binary hosting contract.**

## What Happened

Recorded D033 for the exact ONNX binary hosting contract. Updated the outcome artifact to reference D033 and checked that the decision/outcome preserve the planned-not-uploaded and source-blocked status without leaks or signed URLs.

## Verification

Decision and outcome checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_decision_save D033` | 0 | ✅ pass — decision recorded | 0ms |
| 2 | `gsd_exec M035 decision/outcome checks` | 0 | ✅ pass — required markers present, no leak markers, no signed URLs | 52ms |

## Deviations

None.

## Known Issues

Exact ONNX model binary is not uploaded/mirrored; no hosted source URL exists.

## Files Created/Modified

- `.gsd/DECISIONS.md`
- `benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt`
