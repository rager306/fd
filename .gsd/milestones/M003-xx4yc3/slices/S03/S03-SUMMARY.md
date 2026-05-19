---
id: S03
parent: M003-xx4yc3
milestone: M003-xx4yc3
provides:
  - Validated cache baseline for benchmark interpretation.
requires:
  []
affects:
  - S04
key_files: []
key_decisions:
  - Use Redis key/payload evidence plus API miss-log counts to validate cache tiers until explicit hit metrics exist.
patterns_established:
  - Use `STRLEN` to validate binary embedding payload size: 2 + dim*4 bytes.
observability_surfaces:
  - Redis key scans, STRLEN payload checks, request timing files, and API cache-miss logs.
drill_down_paths:
  - .gsd/milestones/M003-xx4yc3/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S03/tasks/T03-SUMMARY.md
  - .gsd/milestones/M003-xx4yc3/slices/S03/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T08:19:13.818Z
blocker_discovered: false
---

# S03: Runtime cache validation

**S03 proved runtime L1/L2 cache behavior and dimension-scoped Redis storage.**

## What Happened

S03 validated real cache behavior. Redis stores dimension-scoped keys with expected binary payload sizes. Same text requested at 1024d and 512d produced separate Redis keys. Warm requests avoid TEI while L1 is hot. After API restart, L1 is lost but Redis L2 serves the cached key without another TEI miss.

## Verification

All runtime cache validation checks passed.

## Requirements Advanced

- Runtime cache validation complete. — 

## Requirements Validated

- 1024d Redis payload size 4098 bytes. — 
- 512d Redis payload size 2050 bytes. — 
- Same-text dimension-scoped keys exist separately. — 
- Warm request faster than cold and no repeated TEI miss. — 
- Redis L2 serves cached key after API restart. — 

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Warm timing used `/usr/bin/time` with two-decimal precision, so warm path printed 0.00s; log-count evidence confirmed no TEI miss.

## Known Limitations

No explicit L1/L2 hit counters or logs exist; cache hit inference uses timings, Redis state, and absence of miss logs.

## Follow-ups

Run benchmark.py and resource/log correlation in S04.

## Files Created/Modified

None.
