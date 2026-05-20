---
id: T02
parent: S04
milestone: M014-vjfs9f
key_files:
  - .gsd/DECISIONS.md
  - benchmark-results/fd-benchmark-m014-comparison.txt
key_decisions:
  - D010 records the M014 recommendation: continue ONNX experimentally, keep TEI default.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:34:31.828Z
blocker_discovered: false
---

# T02: Recorded the M014 benchmark recommendation as a durable GSD decision.

**Recorded the M014 benchmark recommendation as a durable GSD decision.**

## What Happened

Recorded the benchmark recommendation as GSD decision D010. The decision preserves TEI as the production/default backend and recommends continuing tagged ONNX only as an opt-in experimental path until packaging, repeated performance, and quality gates are satisfied.

## Verification

GSD decision save succeeded and regenerated `.gsd/DECISIONS.md`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_decision_save D010 runtime-performance recommendation` | 0 | ✅ pass — decision saved | 0ms |

## Deviations

None.

## Known Issues

Production switch remains blocked by native artifact packaging, CI/Docker integration, repeated benchmarks, and larger Russian/legal quality evaluation.

## Files Created/Modified

- `.gsd/DECISIONS.md`
- `benchmark-results/fd-benchmark-m014-comparison.txt`
