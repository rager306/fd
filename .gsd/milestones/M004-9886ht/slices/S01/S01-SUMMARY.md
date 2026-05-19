---
id: S01
parent: M004-9886ht
milestone: M004-9886ht
provides:
  - Reliable throughput summary for S02-S03 evidence.
requires:
  []
affects:
  - S02
  - S03
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m004-s01.txt
key_decisions:
  - Keep the S01 fix minimal: store throughput rows and compute max row in summary.
patterns_established:
  - Benchmark summaries should be derived from stored measurement rows rather than loop-local values.
observability_surfaces:
  - benchmark-results/fd-benchmark-m004-s01.txt
drill_down_paths:
  - .gsd/milestones/M004-9886ht/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M004-9886ht/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M004-9886ht/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T10:16:31.892Z
blocker_discovered: false
---

# S01: Benchmark summary correctness

**S01 made benchmark max-throughput summary self-consistent and verified it against live runtime.**

## What Happened

S01 fixed the benchmark summary bug identified in M003. The prior code printed the last throughput loop value and hardcoded 16 concurrent. The updated code records each throughput row and prints the max row by req/s with matching concurrency. Verification under uv Python 3.13 passed, and a parser confirmed the summary matches the table maximum.

## Verification

All S01 tasks complete and verified.

## Requirements Advanced

- Performance benchmark correctness improved. — 

## Requirements Validated

- Summary max throughput equals parsed table max. — 
- GitNexus reports low risk and no affected processes. — 

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Benchmark still prints an imprecise `Cache speedup ... (median text)` label; not fixed in this slice because it is separate from throughput summary correctness.

## Follow-ups

Proceed to S02: cache observability and log-noise reduction.

## Files Created/Modified

- `benchmark.py` — Throughput rows now persist and summary selects max measured row.
- `benchmark-results/fd-benchmark-m004-s01.txt` — Verification benchmark output under uv Python 3.13.
