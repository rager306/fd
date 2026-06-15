---
id: T02
parent: S03
milestone: M047-9fngng
key_files:
  - api/embed/tei.go
  - api/embed/tei_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T08:25:54.296Z
blocker_discovered: false
---

# T02: Added bounded TEI retry and repeated-outage fast-fail behavior.

**Added bounded TEI retry and repeated-outage fast-fail behavior.**

## What Happened

Implemented bounded TEI retry behavior in `TEIClient`. The client now retries retriable network errors and 502/503/504 statuses up to three attempts with context-aware backoff, does not retry context cancellation/deadline, 4xx responses, decode errors, or malformed successful payloads, and tracks consecutive retriable call failures. After three failed calls it opens a short circuit cooldown and subsequent calls fail fast with a `TEI circuit open` error without issuing another HTTP request. Successful calls reset the circuit state. Constructor compatibility is preserved.

## Verification

`cd api && gofmt -w embed/tei.go embed/tei_test.go && go test ./embed` passed with 21 tests. `cd api && go test ./...` passed with 288 tests. Static proof `06c49705-07f9-4c63-add2-85eb6ef673c9` passed for retry/circuit code invariants.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && gofmt -w embed/tei.go embed/tei_test.go && go test ./embed` | 0 | ✅ pass | 32400ms |
| 2 | `cd api && go test ./...` | 0 | ✅ pass | 18000ms |
| 3 | `gsd_exec 06c49705-07f9-4c63-add2-85eb6ef673c9` | 0 | ✅ pass | 131ms |

## Deviations

Used deterministic exponential backoff rather than jitter to keep behavior testable and bounded; defaults are small and context-aware.

## Known Issues

S04 still needs warmup retry/closure and final milestone gates.

## Files Created/Modified

- `api/embed/tei.go`
- `api/embed/tei_test.go`
