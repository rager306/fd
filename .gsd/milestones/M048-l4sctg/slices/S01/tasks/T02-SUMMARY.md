---
id: T02
parent: S01
milestone: M048-l4sctg
key_files:
  - api/cache/hash.go
  - api/cache/tiered.go
  - api/cache/redis.go
  - api/internal/envutil/int.go
  - api/internal/envutil/int_test.go
  - api/fd_v2_cache_integration_test.go
  - api/main.go
  - api/middleware/ratelimit.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T11:02:11.643Z
blocker_discovered: false
---

# T02: Removed dead LRUCache and unified duplicate cache hash/env helpers.

**Removed dead LRUCache and unified duplicate cache hash/env helpers.**

## What Happened

Deleted `api/cache/lru.go` and its dedicated tests. Replaced the only non-self `NewLRUCache` scaffold in `fd_v2_cache_integration_test.go` with a LocalCache-backed test adapter that preserves hit/miss metrics behavior. Added package-local `shortHash` and replaced duplicate `shortCacheKeyHash` and `shortNamespaceHash`. Added `internal/envutil` with `Int` and `PositiveInt`, updated `main` and `ratelimit` to use it, and removed duplicated `getEnvInt`/`envInt` functions.

## Verification

`cd api && go test ./cache` passed with 36 tests. `cd api && go test ./...` passed with 282 tests. Static proof `1453b735-d079-4ce7-9282-08805c13a318` passed for removed LRU files and duplicate helper symbols.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./cache` | 0 | ✅ pass | 9700ms |
| 2 | `cd api && go test ./...` | 0 | ✅ pass | 9600ms |
| 3 | `gsd_exec 1453b735-d079-4ce7-9282-08805c13a318` | 0 | ✅ pass | 144ms |

## Deviations

The LocalCache-backed integration test adapter records cache hit/miss metrics to preserve the original fd_v2 cache integration assertion that had been supplied by LRU metrics.

## Known Issues

S02/S03 issue #7 findings remain open.

## Files Created/Modified

- `api/cache/hash.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`
- `api/internal/envutil/int.go`
- `api/internal/envutil/int_test.go`
- `api/fd_v2_cache_integration_test.go`
- `api/main.go`
- `api/middleware/ratelimit.go`
