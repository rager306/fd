---
id: T02
parent: S02
milestone: M047-9fngng
key_files:
  - api/main.go
  - api/main_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T08:17:38.780Z
blocker_discovered: false
---

# T02: Replaced listener goroutine os.Exit with controlled shutdown signalling.

**Replaced listener goroutine os.Exit with controlled shutdown signalling.**

## What Happened

Added `serverErrorSignal` and `reportHTTPServerError`. The helper logs listener startup, ignores wrapped `http.ErrServerClosed` with `errors.Is`, logs non-server-closed listener failures, and sends a synthetic `server_error` signal into the same signal channel used by lifecycle graceful shutdown. `main` now creates the signal channel before starting the listener goroutine and calls `go reportHTTPServerError(...)` instead of exiting inside the goroutine.

## Verification

`cd api && gofmt -w main.go main_test.go && go test ./...` passed with 285 tests in 9 packages. Static proof `519aee78-cfa7-47d0-9fdf-aee5cddd1f83` verified `errors.Is`, no listener-helper `os.Exit`, and routing through `reportHTTPServerError`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && gofmt -w main.go main_test.go && go test ./...` | 0 | ✅ pass | 11700ms |
| 2 | `gsd_exec 519aee78-cfa7-47d0-9fdf-aee5cddd1f83` | 0 | ✅ pass | 108ms |

## Deviations

Used a synthetic `os.Signal` implementation to reuse the existing `lifecycle.AwaitSignalAndShutdown` path with minimal churn.

## Known Issues

S03/S04 still need TEI retry/fast-fail and warmup retry.

## Files Created/Modified

- `api/main.go`
- `api/main_test.go`
