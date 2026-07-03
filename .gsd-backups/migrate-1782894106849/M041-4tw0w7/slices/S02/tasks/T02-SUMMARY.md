---
id: T02
parent: S02
milestone: M041-4tw0w7
key_files:
  - api/lifecycle/warmup.go
  - api/lifecycle/warmup_test.go
  - api/main.go
  - api/main_test.go
  - benchmark-results/m041-s02-t02-go-test.txt
  - benchmark-results/m041-s02-t02-lint.txt
  - benchmark-results/m041-s02-t02-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T05:44:44.513Z
blocker_discovered: false
---

# T02: Added async model pre-warm: server startup no longer blocks on warmup, lifecycle state flips ready only after successful dummy embedding.

**Added async model pre-warm: server startup no longer blocks on warmup, lifecycle state flips ready only after successful dummy embedding.**

## What Happened

Implemented `api/lifecycle/warmup.go` with `PreWarm(ctx, model)` using the existing Embedder shape (`Embed(ctx, []string)`) and one dummy warmup input. It validates exactly one non-empty embedding and wraps model/context errors. Added warmup unit tests for success input, wrapped model error, malformed responses, context cancellation, and nil model. Wired `main.go` with `defaultWarmupTimeout=30s`, `lifecycle.DefaultState()`, and `startModelWarmup`: after the HTTP server goroutine is started, warmup runs asynchronously, logs `model warmup started`, marks lifecycle ready on success, and records `lastError` + logs `model warmup failed` on failure. Added main tests for successful warmup marking readiness and failed warmup storing last error while staying not-ready.

## Verification

Fresh verification passed: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. Evidence files: benchmark-results/m041-s02-t02-go-test.txt, benchmark-results/m041-s02-t02-lint.txt, benchmark-results/m041-s02-t02-govulncheck.txt. GitNexus detect_changes reports MEDIUM risk because api/main.go startup flow changed, no high/critical risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

The original T02 verification text mentions `/live` and `/ready` smoke checks. Those endpoints are explicitly planned in T03 and do not exist yet, so endpoint-level startup smoke is deferred to T03/T06. T02 is verified at primitive/main-helper level plus full M043 gate.

## Known Issues

`api/report.json` remains an unrelated untracked generated file outside this task. T02 changes are not committed yet.

## Files Created/Modified

- `api/lifecycle/warmup.go`
- `api/lifecycle/warmup_test.go`
- `api/main.go`
- `api/main_test.go`
- `benchmark-results/m041-s02-t02-go-test.txt`
- `benchmark-results/m041-s02-t02-lint.txt`
- `benchmark-results/m041-s02-t02-govulncheck.txt`
