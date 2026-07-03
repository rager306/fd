---
id: T02
parent: S05
milestone: M046-zqzcu6
key_files:
  - api/cache/local.go
  - api/cache/local_test.go
  - api/main.go
key_decisions:
  - Use a single mutex-owned map and derived length instead of trying to repair a separate counter around sync.Map.
duration: 
verification_result: passed
completed_at: 2026-06-14T19:18:33.242Z
blocker_discovered: false
---

# T02: Refactored LocalCache to a single mutex-owned map with idempotent Close.

**Refactored LocalCache to a single mutex-owned map with idempotent Close.**

## What Happened

Replaced `sync.Map` plus separate size counter with one mutex-owned `map[string]l1Entry`, so size is derived from `len(data)` and cannot drift independently. Added `Close() error` with `sync.Once` to stop the background eviction loop. `Get`, `Set`, `Delete`, max-size enforcement, lazy expiry, and background expiry behavior remain compatible. Added tests for close idempotency and concurrent overwrite accounting.

## Verification

Targeted LocalCache tests passed and `cd api && go test -race ./cache -run TestLocalCache` passed. `cd api && go test ./cache && go test ./...` passed with 44 cache tests and 281 total tests. Static proof `f124000a-5996-4c68-888d-1e31237c6d39` passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./cache -run 'TestLocalCache_(CloseIsIdempotent|ConcurrentOverwriteKeepsSingleEntry|SetAndGet|TTLExpired|SetRefreshesExistingValueAndTTL|EnforcesMaxSize|Delete|NotFound)'` | 0 | ✅ pass | 35000ms |
| 2 | `cd api && go test -race ./cache -run TestLocalCache` | 0 | ✅ pass | 10000ms |
| 3 | `cd api && go test ./cache && go test ./...` | 0 | ✅ pass | 10000ms |
| 4 | `gsd_exec f124000a-5996-4c68-888d-1e31237c6d39` | 0 | ✅ pass | 91ms |

## Deviations

None.

## Known Issues

Quality gates and artifact/UAT remain in T03/T04.

## Files Created/Modified

- `api/cache/local.go`
- `api/cache/local_test.go`
- `api/main.go`
