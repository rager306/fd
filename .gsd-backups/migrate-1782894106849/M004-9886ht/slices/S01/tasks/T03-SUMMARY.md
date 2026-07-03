---
id: T03
parent: S01
milestone: M004-9886ht
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m004-s01.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T10:16:13.920Z
blocker_discovered: false
---

# T03: Verified benchmark throughput summary now matches the measured table max.

**Verified benchmark throughput summary now matches the measured table max.**

## What Happened

Ran the fixed benchmark under uv Python 3.13 against the live Docker stack and saved the output. A parser confirmed the throughput table maximum was 742.8 req/s at concurrency 16 and the summary reported ~743 req/s at 16 concurrent. GitNexus change detection after the edit remained low risk with no affected processes.

## Verification

`uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m004-s01.txt` passed. Parser confirmed table max and summary match. GitNexus change detection low risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m004-s01.txt` | 0 | ✅ pass: summary reports ~743 req/s at 16 concurrent | 22800ms |
| 2 | `parse benchmark-results/fd-benchmark-m004-s01.txt and compare throughput table max to summary` | 0 | ✅ pass: table_max=742.8@16; summary=743@16 | 0ms |
| 3 | `gitnexus_detect_changes(repo=fd, scope=all)` | 0 | ✅ pass: low risk; no affected processes | 0ms |

## Deviations

The verification run naturally produced different throughput values from M003; in this run the true maximum was again at 16 concurrency, so the important proof is the parsed equality between table max and summary, not a fixed concurrency value.

## Known Issues

The existing `Cache speedup: ... (median text)` label remains imprecise because it reports max speedup, but it is outside this slice's throughput-summary fix.

## Files Created/Modified

- `benchmark.py`
- `benchmark-results/fd-benchmark-m004-s01.txt`
