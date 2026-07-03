---
id: T02
parent: S03
milestone: M049-7dn2gp
key_files:
  - benchmark-results/m049-s03-live-container-proof.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T13:13:57.372Z
blocker_discovered: false
---

# T02: Rebuilt the API container and passed live health/metrics/cache invalidation smoke.

**Rebuilt the API container and passed live health/metrics/cache invalidation smoke.**

## What Happened

Rebuilt/restarted the `api` service with Docker Compose, waited for `/health`, and ran authenticated live HTTP checks against the rebuilt container. The smoke proved new health dependency/capacity fields, metrics gauges, auth protection on cache flush, MISS->HIT->flush->MISS behavior, and MISS->HIT->delete->MISS behavior. Docker Compose reports api, redis, and tei healthy after the run.

## Verification

`docker compose up -d --build api` completed. `/health` returned ok with `in_flight_capacity`, TEI dependency reachable, and Redis dependency reachable. Live smoke artifact `benchmark-results/m049-s03-live-container-proof.md` reports `SUMMARY passed=5 failed=0 total=5`. `docker compose ps api redis tei` shows all three services healthy.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker compose up -d --build api` | 0 | ✅ pass | 27400ms |
| 2 | `M049 live container smoke` | 0 | ✅ pass | 10000ms |
| 3 | `docker compose ps api redis tei` | 0 | ✅ pass | 0ms |

## Deviations

None.

## Known Issues

None for runtime verification. Closure artifacts/milestone validation remain.

## Files Created/Modified

- `benchmark-results/m049-s03-live-container-proof.md`
