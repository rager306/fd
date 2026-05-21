---
id: T01
parent: S02
milestone: M036-o0hewj
key_files:
  - .gsd/DECISIONS.md
  - benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt
key_decisions:
  - D034: planned reproducible-export workflow contract is the no-upload alternative, not current regenerated-export proof.
duration: 
verification_result: passed
completed_at: 2026-05-21T09:38:45.377Z
blocker_discovered: false
---

# T01: Recorded D034 for the planned reproducible-export workflow contract.

**Recorded D034 for the planned reproducible-export workflow contract.**

## What Happened

Recorded D034 and updated the outcome artifact to reference it. Decision/outcome checks confirm the planned_not_proven boundary and the M032 verifier claim scope remain explicit, with no leak markers or signed URLs.

## Verification

Decision and outcome checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_decision_save D034` | 0 | ✅ pass — decision recorded | 0ms |
| 2 | `gsd_exec M036 decision outcome checks` | 0 | ✅ pass — required markers present, no leak markers, no signed URLs | 46ms |

## Deviations

None.

## Known Issues

No regenerated export proof exists yet.

## Files Created/Modified

- `.gsd/DECISIONS.md`
- `benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt`
