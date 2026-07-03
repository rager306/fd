---
id: T01
parent: S01
milestone: M049-7dn2gp
key_files:
  - api/cache/local_test.go
  - api/cache/redis_test.go
  - api/cache/tiered_test.go
  - api/handlers/cache_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T12:50:01.180Z
blocker_discovered: false
---

# T01: Pinned cache invalidation gaps with red tests.

**Pinned cache invalidation gaps with red tests.**

## What Happened

Added red tests for LocalCache Flush/Size, TieredCache Delete/Flush/LocalSize, Redis namespace pattern safety, and a new HTTP cache invalidation handler contract. The tests fail at compile time because the invalidation primitives and handler do not exist yet, which matches issue #8 AN-A.

## Verification

`cd api && go test ./cache ./handlers` failed as expected: missing `LocalCache.Flush`, `LocalCache.Size`, `RedisCache.namespacePattern`, `TieredCache.Delete`, `TieredCache.Flush`, `TieredCache.LocalSize`, `CacheHandler`, and `NewCacheHandler`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./cache ./handlers` | 1 | ✅ expected red | 10000ms |

## Deviations

None.

## Known Issues

S01 remains red until cache primitives and handler are implemented.

## Files Created/Modified

- `api/cache/local_test.go`
- `api/cache/redis_test.go`
- `api/cache/tiered_test.go`
- `api/handlers/cache_test.go`
