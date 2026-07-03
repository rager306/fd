---
id: T01
parent: S01
milestone: M004-9886ht
key_files:
  - benchmark.py
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T10:14:48.302Z
blocker_discovered: false
---

# T01: Identified throughput summary bug: summary uses final loop value instead of max throughput row.

**Identified throughput summary bug: summary uses final loop value instead of max throughput row.**

## What Happened

Inspected benchmark.py and confirmed the root cause of the M003 throughput summary discrepancy. The throughput loop computes `rps` per concurrency but does not persist rows. The final summary then prints the last `rps` value and hardcodes `16 concurrent`, so it reports the final loop iteration rather than the maximum measured throughput row. GitNexus impact analysis for `Function:benchmark.py:main` returned LOW risk with one direct file-level impact and no affected processes.

## Verification

GitNexus impact analysis ran for `Function:benchmark.py:main` before editing: LOW risk, no affected processes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gitnexus_impact(target='Function:benchmark.py:main', direction='upstream', repo='fd')` | 0 | ✅ pass: LOW risk; direct impact only benchmark.py; no affected processes | 0ms |

## Deviations

Python LSP is unavailable in this environment, so benchmark.py was inspected with file read and GitNexus impact analysis by symbol UID.

## Known Issues

None.

## Files Created/Modified

- `benchmark.py`
