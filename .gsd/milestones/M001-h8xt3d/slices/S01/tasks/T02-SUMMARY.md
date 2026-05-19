---
id: T02
parent: S01
milestone: M001-h8xt3d
key_files:
  - api/cache/redis.go
  - api/cache/tiered.go
  - api/cache/redis_binary_test.go
  - api/cache/tiered_cache_test.go
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T06:51:15.831Z
blocker_discovered: false
---

# T02: Fixed cache dimension isolation and short-vector panic risk.

**Fixed cache dimension isolation and short-vector panic risk.**

## What Happened

Implemented dimension-aware L1 and singleflight keys in TieredCache so the same text requested at 512d and 1024d no longer shares local cache or in-flight results. Changed marshalEmbedding to validate dimensions and short vectors, returning errors instead of panicking. RedisCache.Set now propagates marshal validation errors, and TieredCache validates embeddings before L1/L2 backfill. Added focused cache tests for 512d/1024d isolation and short-vector error behavior.

## Verification

Targeted cache tests passed: `cd api && go test ./cache -run 'Test.*(Tiered|Binary|Redis|Local|Marshal)' -count=1` reported 13 tests passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./cache -run 'Test.*(Tiered|Binary|Redis|Local|Marshal)' -count=1` | 0 | ✅ pass: 13 tests passed in 1 package | 4200ms |

## Deviations

Used gofmt because no Go LSP formatter is available in this environment.

## Known Issues

None.

## Files Created/Modified

- `api/cache/redis.go`
- `api/cache/tiered.go`
- `api/cache/redis_binary_test.go`
- `api/cache/tiered_cache_test.go`
