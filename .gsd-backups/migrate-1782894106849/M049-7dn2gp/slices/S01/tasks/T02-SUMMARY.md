---
id: T02
parent: S01
milestone: M049-7dn2gp
key_files:
  - api/cache/local.go
  - api/cache/redis.go
  - api/cache/tiered.go
  - api/handlers/cache.go
  - api/main.go
  - api/cache/local_test.go
  - api/cache/redis_test.go
  - api/cache/tiered_test.go
  - api/handlers/cache_test.go
key_decisions:
  - Use `POST /v1/cache/delete` with input+dimensions instead of `DELETE /v1/cache/:keyHash` because key hashes are not reversible and input+dimension matches the actual key derivation.
duration: 
verification_result: passed
completed_at: 2026-06-15T12:52:41.283Z
blocker_discovered: false
---

# T02: Implemented namespace-safe cache invalidation primitives and HTTP routes.

**Implemented namespace-safe cache invalidation primitives and HTTP routes.**

## What Happened

Added `LocalCache.Flush`/`Size`, `RedisCache.Delete`/`FlushNamespace`, and `TieredCache.Delete`/`Flush`/`LocalSize`. Redis flush scans and deletes only keys matching the configured fd cache namespace; it never calls `FlushDB`. Added `handlers.CacheHandler` with `POST /v1/cache/flush` and `POST /v1/cache/delete`, wired through the existing global API-key middleware in `main.go`. Delete accepts string or string-array input and defaults dimensions to 1024 while validating the existing 512/1024 dimensions.

## Verification

`cd api && go test ./cache ./handlers` passed with 127 tests. `cd api && go test ./...` passed with 293 tests.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./cache ./handlers` | 0 | ✅ pass | 5000ms |
| 2 | `cd api && go test ./...` | 0 | ✅ pass | 7200ms |

## Deviations

The HTTP delete route uses input text plus dimensions instead of `:keyHash`; current cache keys are derived from text+dimension and short hashes are not reversible identifiers. This is safer and directly actionable for a solo operator.

## Known Issues

Live container cache HIT->flush->MISS proof is deferred to S03 runtime verification.

## Files Created/Modified

- `api/cache/local.go`
- `api/cache/redis.go`
- `api/cache/tiered.go`
- `api/handlers/cache.go`
- `api/main.go`
- `api/cache/local_test.go`
- `api/cache/redis_test.go`
- `api/cache/tiered_test.go`
- `api/handlers/cache_test.go`
