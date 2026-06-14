---
id: T05
parent: S02
milestone: M041-4tw0w7
key_files:
  - api/lifecycle/shutdown.go
  - api/lifecycle/shutdown_test.go
  - api/main.go
  - benchmark-results/m041-s02-t05-go-test.txt
  - benchmark-results/m041-s02-t05-lint.txt
  - benchmark-results/m041-s02-t05-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T06:11:43.201Z
blocker_discovered: false
---

# T05: Added lifecycle-managed graceful shutdown for SIGTERM/SIGINT with 30s shared server/in-flight drain and force-close failure path.

**Added lifecycle-managed graceful shutdown for SIGTERM/SIGINT with 30s shared server/in-flight drain and force-close failure path.**

## What Happened

Implemented `api/lifecycle/shutdown.go` with `AwaitSignalAndShutdown` and `GracefulShutdown`. On SIGTERM/SIGINT the orchestrator marks lifecycle shutdown immediately, logs the signal and timeout, runs `http.Server.Shutdown(ctx)` and `State.WaitDrain(timeout)` under one shared deadline, and force-closes the server with a non-nil error if the deadline is exceeded. Updated `main.go` to delegate signal handling to lifecycle, use `lifecycle.DefaultShutdownTimeout` (30s), exit 1 on failed/forced shutdown, and close ONNX/Redis resources explicitly so `os.Exit(1)` does not bypass defers. Added lifecycle shutdown tests for idle clean exit under one second, state marked shutting down before server shutdown, in-flight drain completion, and force-close on timeout.

## Verification

Fresh M043 gate passed: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. Evidence files: benchmark-results/m041-s02-t05-go-test.txt, benchmark-results/m041-s02-t05-lint.txt, benchmark-results/m041-s02-t05-govulncheck.txt. GitNexus detect_changes reports MEDIUM risk because `api/main.go` shutdown flow changed, no high/critical risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

The T05 plan requested a full process-level SIGTERM integration test with a long-running request. This task adds deterministic lifecycle/server unit coverage for the shutdown orchestrator; the full runtime process integration is deferred to T06, which is explicitly scoped for integration tests across F-1/F-2/F-5 and now has all lifecycle pieces in place.

## Known Issues

`api/report.json` remains an unrelated untracked generated file outside this task. T05 changes are not committed yet.

## Files Created/Modified

- `api/lifecycle/shutdown.go`
- `api/lifecycle/shutdown_test.go`
- `api/main.go`
- `benchmark-results/m041-s02-t05-go-test.txt`
- `benchmark-results/m041-s02-t05-lint.txt`
- `benchmark-results/m041-s02-t05-govulncheck.txt`
