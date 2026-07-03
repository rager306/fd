---
id: T01
parent: S05
milestone: M046-zqzcu6
key_files:
  - api/cache/local_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T19:16:38.315Z
blocker_discovered: false
---

# T01: Added red tests for LocalCache close lifecycle and concurrent overwrite accounting.

**Added red tests for LocalCache close lifecycle and concurrent overwrite accounting.**

## What Happened

Added tests requiring `LocalCache.Close()` to exist and be idempotent, and requiring concurrent overwrites of one key to count as a single retained entry. The targeted test run failed at compile time because `LocalCache` has no `Close` method, establishing the lifecycle gap before implementation.

## Verification

`cd api && go test ./cache -run 'TestLocalCache_(CloseIsIdempotent|ConcurrentOverwriteKeepsSingleEntry)'` failed with `c.Close undefined`, as expected before implementation.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./cache -run 'TestLocalCache_(CloseIsIdempotent|ConcurrentOverwriteKeepsSingleEntry)'` | 1 | ✅ expected red fail | 1000ms |

## Deviations

The first red state is compile-red for the new lifecycle surface. The concurrency/accounting behavior will be exercised after the package compiles.

## Known Issues

Implementation pending in T02.

## Files Created/Modified

- `api/cache/local_test.go`
