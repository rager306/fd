---
id: T01
parent: S04
milestone: M047-9fngng
key_files:
  - api/main_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T08:29:57.708Z
blocker_discovered: false
---

# T01: Pinned warmup retry behavior with red tests.

**Pinned warmup retry behavior with red tests.**

## What Happened

Added tests for an injectable warmup retry policy. The red run proves `startModelWarmupWithPolicy` and `warmupRetryPolicy` do not yet exist, covering the desired contract: retry after transient warmup failure and mark ready on later success, and record terminal error after bounded attempts without marking ready.

## Verification

`cd api && go test ./...` failed as expected with undefined `startModelWarmupWithPolicy` and `warmupRetryPolicy`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./...` | 1 | ✅ expected red | 8600ms |

## Deviations

None.

## Known Issues

S04 remains red until bounded warmup retry implementation is added.

## Files Created/Modified

- `api/main_test.go`
