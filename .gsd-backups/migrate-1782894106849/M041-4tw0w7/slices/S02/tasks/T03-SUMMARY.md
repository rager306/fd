---
id: T03
parent: S02
milestone: M041-4tw0w7
key_files:
  - api/handlers/probes.go
  - api/handlers/probes_test.go
  - api/handlers/health.go
  - api/main.go
  - benchmark-results/m041-s02-t03-go-test.txt
  - benchmark-results/m041-s02-t03-lint.txt
  - benchmark-results/m041-s02-t03-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T05:53:31.178Z
blocker_discovered: false
---

# T03: Added /live and /ready probes backed by lifecycle state: /live always 200, /ready 200 after warmup and 503 model_not_loaded before warmup/shutdown.

**Added /live and /ready probes backed by lifecycle state: /live always 200, /ready 200 after warmup and 503 model_not_loaded before warmup/shutdown.**

## What Happened

Implemented `api/handlers/probes.go` with `NewLiveHandler` and `NewReadyHandler`. `/live` returns cheap 200 with status/time and does not depend on model readiness. `/ready` uses `lifecycle.State.IsReady()`: returns 200 status=ready after warmup, otherwise sets `Retry-After: 5` and emits the existing OpenAI-style `model_not_loaded` overloaded_error envelope. Wired `main.go` to register `GET /live` and `GET /ready` using the lifecycle state created for warmup. Added `api/handlers/probes_test.go` covering live OK, ready 503 before warmup with code/type/header, ready 200 after MarkWarmupDone, and ready 503 during shutdown. goconst flagged duplicated `status`/`time` keys across health/probes; fixed by introducing shared unexported constants in `handlers/health.go`.

## Verification

Fresh M043 gate passed: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. Evidence files: benchmark-results/m041-s02-t03-go-test.txt, benchmark-results/m041-s02-t03-lint.txt, benchmark-results/m041-s02-t03-govulncheck.txt. GitNexus detect_changes reports MEDIUM risk due health/main route changes, no high/critical risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

Endpoint-level tests are unit-level gin router tests, not a live binary curl smoke. The live startup sequence (`/live` immediately 200 and `/ready` 503 then 200 after async warmup) is planned for later S02 integration tests (T06) once lifecycle gate/shutdown wiring is complete.

## Known Issues

`api/report.json` remains an unrelated untracked generated file outside this task. T03 changes are not committed yet.

## Files Created/Modified

- `api/handlers/probes.go`
- `api/handlers/probes_test.go`
- `api/handlers/health.go`
- `api/main.go`
- `benchmark-results/m041-s02-t03-go-test.txt`
- `benchmark-results/m041-s02-t03-lint.txt`
- `benchmark-results/m041-s02-t03-govulncheck.txt`
