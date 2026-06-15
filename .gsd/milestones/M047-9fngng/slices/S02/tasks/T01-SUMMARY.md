---
id: T01
parent: S02
milestone: M047-9fngng
key_files:
  - api/main_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T08:16:08.867Z
blocker_discovered: false
---

# T01: Pinned listener fatal-error routing with red tests.

**Pinned listener fatal-error routing with red tests.**

## What Happened

Added tests for a new listener error reporting helper. The red run proves the helper and `serverErrorSignal` type do not yet exist, covering the intended contract: wrapped `http.ErrServerClosed` is ignored via `errors.Is`, while arbitrary listener errors produce a synthetic shutdown signal for main control flow.

## Verification

`cd api && go test ./...` failed as expected with undefined `reportHTTPServerError` and `serverErrorSignal`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./...` | 1 | ✅ expected red | 9000ms |

## Deviations

None.

## Known Issues

S02 remains red until the listener helper and main control-flow wiring are implemented.

## Files Created/Modified

- `api/main_test.go`
