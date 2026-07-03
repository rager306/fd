---
id: T03
parent: S03
milestone: M004-9886ht
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m004-s03.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T10:29:49.299Z
blocker_discovered: false
---

# T03: Verified benchmark now reports Redis L2 after-restart behavior under uv Python 3.13.

**Verified benchmark now reports Redis L2 after-restart behavior under uv Python 3.13.**

## What Happened

Ran the updated benchmark under uv Python 3.13 and saved the output. The new section primed a dedicated text, restarted the API through Docker Compose, waited for health, then requested the same text. Result: prime/cold request 299.11ms; after API restart 3.10ms, demonstrating Redis L2 persistence across API process restart. The throughput summary also remained self-consistent: table max 754.2 req/s at concurrency 4 and summary ~754 req/s at concurrency 4. Compose status after the run showed all services healthy.

## Verification

`uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m004-s03.txt` passed. Parser confirmed throughput summary consistency and Redis L2 result. GitNexus change detection stayed low risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests --with redis python benchmark.py | tee benchmark-results/fd-benchmark-m004-s03.txt` | 0 | ✅ pass: Redis L2 after restart 3.10ms | 25600ms |
| 2 | `parse benchmark-results/fd-benchmark-m004-s03.txt` | 0 | ✅ pass: table_max=754.2@4; summary=754@4; l2_restart=3.10ms | 0ms |
| 3 | `docker compose ps` | 0 | ✅ pass: api, redis, tei healthy | 0ms |
| 4 | `gitnexus_detect_changes(repo=fd, scope=all)` | 0 | ✅ pass: low risk; no affected processes | 0ms |

## Deviations

The benchmark's Redis L2 diagnostic restarts the API container during the run; this is expected and verified the API returned healthy afterward.

## Known Issues

Cold latency varied upward during this run, but the diagnostic target passed: after API restart, the primed text returned in 3.10ms from Redis L2/backfilled L1.

## Files Created/Modified

- `benchmark.py`
- `benchmark-results/fd-benchmark-m004-s03.txt`
