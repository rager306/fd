---
id: T04
parent: S02
milestone: M041-4tw0w7
key_files:
  - api/middleware/lifecycle.go
  - api/middleware/lifecycle_test.go
  - api/main.go
  - benchmark-results/m041-s02-t04-go-test.txt
  - benchmark-results/m041-s02-t04-lint.txt
  - benchmark-results/m041-s02-t04-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T06:07:30.971Z
blocker_discovered: false
---

# T04: Added lifecycle gate middleware for /v1/embeddings with warmup/shutdown rejection and in-flight request tracking.

**Added lifecycle gate middleware for /v1/embeddings with warmup/shutdown rejection and in-flight request tracking.**

## What Happened

Implemented `api/middleware/lifecycle.go` with `LifecycleGate(state)`. The middleware rejects requests during shutdown first (`shutting_down`, `Retry-After: 30`), rejects unready/warmup state (`model_not_loaded`, `Retry-After: 5`), and tracks accepted requests with `state.TrackRequest()` until downstream handlers finish. Wired `/v1/embeddings` in `main.go` as validation -> lifecycle gate -> embedding handler, preserving the existing validation-before-inference behavior. Added `api/middleware/lifecycle_test.go` for pre-warmup rejection, shutdown rejection, and in-flight tracking using `WaitDrain(0)` while a downstream handler is blocked.

## Verification

Fresh M043 gate passed: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. Evidence files: benchmark-results/m041-s02-t04-go-test.txt, benchmark-results/m041-s02-t04-lint.txt, benchmark-results/m041-s02-t04-govulncheck.txt. GitNexus detect_changes reports LOW risk for currently indexed changed symbols; untracked new middleware symbols are covered by tests/lint before staging.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

T04 gates only `/v1/embeddings` as planned. `/embeddings/batch` is not changed in this task because the T04 plan explicitly scopes the lifecycle gate to `/v1/embeddings`.

## Known Issues

`api/report.json` remains an unrelated untracked generated file outside this task. T04 changes are not committed yet.

## Files Created/Modified

- `api/middleware/lifecycle.go`
- `api/middleware/lifecycle_test.go`
- `api/main.go`
- `benchmark-results/m041-s02-t04-go-test.txt`
- `benchmark-results/m041-s02-t04-lint.txt`
- `benchmark-results/m041-s02-t04-govulncheck.txt`
