---
id: S01
parent: M048-l4sctg
milestone: M048-l4sctg
provides:
  - R037 validated.
  - Issue #7 findings #19, #27, and #28 closed.
requires:
  []
affects:
  []
key_files:
  - api/cache/hash.go
  - api/cache/tiered.go
  - api/cache/redis.go
  - api/internal/envutil/int.go
  - api/internal/envutil/int_test.go
  - api/fd_v2_cache_integration_test.go
  - api/main.go
  - api/middleware/ratelimit.go
  - benchmark-results/m048-s01-cache-cleanup.md
key_decisions:
  - Delete LRUCache rather than marking it reserved because current production code uses LocalCache/TieredCache.
  - Use envutil.Int and envutil.PositiveInt to preserve main zero-allowed vs rate-limit positive-only semantics.
patterns_established:
  - Dead test scaffolds should be replaced with active-path adapters before deleting unused production code.
observability_surfaces:
  - S01 evidence artifact records pre-fix and post-fix static proof plus test results.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-15T11:05:02.471Z
blocker_discovered: false
---

# S01: Cache cleanup consolidation

**Removed dead LRUCache and unified duplicate cache hash/env parsing helpers.**

## What Happened

S01 resolved issue #7 findings #19, #27, and #28. Pre-fix static proof confirmed LRUCache existed with only a test scaffold reference, duplicate short hash helpers existed, and env integer parsers were duplicated. The implementation deleted LRUCache source/tests, replaced the fd_v2 cache integration scaffold with a LocalCache-backed adapter that preserves HIT/MISS metrics assertions, added canonical cache `shortHash`, and introduced `internal/envutil` for active integer configuration parsing. Main and rate-limit middleware now use envutil; duplicate parser functions are gone. R037 was validated.

## Verification

Pre-fix proof `12ffe3b3-84f7-4e6f-8f6a-e3a12f9eef57` passed. `go test ./cache` passed with 36 tests. `go test ./...` passed with 282 tests. Post-cleanup static proof `1453b735-d079-4ce7-9282-08805c13a318` passed. Artifact completeness `52d98836-5c63-4aab-b9ba-72377f58ba41` passed. UAT PASS saved with evidence `2ae5e91d-6c8a-48f7-9e82-505921af6680`, `e7475039-5ae6-4261-b8cb-b3e48ad50841`, `1453b735-d079-4ce7-9282-08805c13a318`, and `c26cc783-387c-41dd-8dbe-13c521b29e34`.

## Requirements Advanced

None.

## Requirements Validated

- R037 — S01 tests and static proof validate LRU removal, hash helper unification, and envutil consolidation.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

UAT command guard blocked a direct source scan containing environment parsing keywords; S01 cites the prior approved gsd_exec static proof for that check.

## Known Limitations

S02/S03 issue #7 findings remain open.

## Follow-ups

Proceed to S02 runtime contract simplification.

## Files Created/Modified

- `api/cache/lru.go` — Deleted dead LRU cache production code.
- `api/cache/lru_test.go` — Deleted dedicated tests for removed LRU cache.
- `api/cache/lru_rapid_test.go` — Deleted dedicated rapid tests for removed LRU cache.
- `api/cache/hash.go` — Added canonical shortHash helper.
- `api/cache/tiered.go` — Uses shortHash.
- `api/cache/redis.go` — Uses shortHash.
- `api/internal/envutil/int.go` — Added shared integer configuration parsing helpers.
- `api/fd_v2_cache_integration_test.go` — Uses LocalCache-backed test adapter instead of LRU.
- `benchmark-results/m048-s01-cache-cleanup.md` — S01 evidence artifact.
