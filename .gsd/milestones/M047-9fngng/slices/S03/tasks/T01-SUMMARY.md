---
id: T01
parent: S03
milestone: M047-9fngng
key_files:
  - api/embed/tei_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T08:22:42.999Z
blocker_discovered: false
---

# T01: Pinned TEI retry and fast-fail gaps with red tests.

**Pinned TEI retry and fast-fail gaps with red tests.**

## What Happened

Added TEI client tests proving retriable 503 should retry and succeed on a later attempt, 400 should not retry, and repeated 503 failures should open a circuit so the next call fails without another HTTP request. The current implementation fails the retry and circuit tests, confirming issue #6 #11.

## Verification

`cd api && go test ./embed` failed as expected: `TestTEIClientEmbedRetriesRetriableStatus` returned status 503 without retry, and `TestTEIClientEmbedFastFailsAfterRepeatedRetriableFailures` returned status 503 instead of `TEI circuit open`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./embed` | 1 | ✅ expected red | 10100ms |

## Deviations

None.

## Known Issues

S03 remains red until retry and fast-fail behavior are implemented.

## Files Created/Modified

- `api/embed/tei_test.go`
