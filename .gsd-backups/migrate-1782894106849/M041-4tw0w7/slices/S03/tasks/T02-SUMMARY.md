---
id: T02
parent: S03
milestone: M041-4tw0w7
key_files:
  - api/handlers/observability.go
  - api/handlers/observability_test.go
  - api/lifecycle/state.go
  - api/lifecycle/state_test.go
  - api/main.go
  - api/handlers/health_test.go
  - benchmark-results/m041-s03-t02-go-test.txt
  - benchmark-results/m041-s03-t02-lint.txt
  - benchmark-results/m041-s03-t02-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T06:34:22.315Z
blocker_discovered: false
---

# T02: Added /version, /info, and /v1/healthcheck observability endpoints backed by buildinfo and lifecycle state.

**Added /version, /info, and /v1/healthcheck observability endpoints backed by buildinfo and lifecycle state.**

## What Happened

Implemented `api/handlers/observability.go` with `NewVersionHandler` and `NewInfoHandler`. `/version` returns build metadata (service, version, model, build hash/date, started_at) plus uptime string and uptime_seconds. `/info` returns model metadata with dims `[512,1024]`, max input tokens 512, max batch size 32, device `cpu`, and lifecycle-derived `loaded`/`warmup_done`. Added `State.IsWarmupDone()` so observability can distinguish successful warmup from readiness, which may be false during shutdown. Wired `main.go` to create one process `buildinfo.Info`, expose `/version`, `/info`, and mount `/v1/healthcheck` as the same handler as `/health`. Added unit tests for version response fields/uptime, info lifecycle state before/after warmup, and healthcheck alias behavior.

## Verification

Fresh M043 gate passed: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. Evidence files: benchmark-results/m041-s03-t02-go-test.txt, benchmark-results/m041-s03-t02-lint.txt, benchmark-results/m041-s03-t02-govulncheck.txt. Targeted tests for /version, /info, /v1/healthcheck, and lifecycle warmup state also passed. GitNexus detect_changes reports MEDIUM risk due `api/main.go` route wiring, no high/critical risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

T02 verification text mentions `/metrics` and deep `/health`; those are explicitly owned by later S03 tasks T03 and T04. This task implemented only the T02 endpoint scope: `/version`, `/info`, and `/v1/healthcheck`. `device` defaults to `cpu` because no runtime provider/device metadata exists yet in RuntimeHealth; this can be enriched by later ONNX/TEI observability work.

## Known Issues

`api/report.json` remains an unrelated untracked generated file outside this task. T02 changes are not committed yet.

## Files Created/Modified

- `api/handlers/observability.go`
- `api/handlers/observability_test.go`
- `api/lifecycle/state.go`
- `api/lifecycle/state_test.go`
- `api/main.go`
- `api/handlers/health_test.go`
- `benchmark-results/m041-s03-t02-go-test.txt`
- `benchmark-results/m041-s03-t02-lint.txt`
- `benchmark-results/m041-s03-t02-govulncheck.txt`
