---
id: T06
parent: S02
milestone: M041-4tw0w7
key_files:
  - api/fd_v2_lifecycle_integration_test.go
  - api/lifecycle/state.go
  - api/lifecycle/state_test.go
  - api/middleware/lifecycle.go
  - api/middleware/lifecycle_test.go
  - api/main.go
  - benchmark-results/m041-s02-t06-lifecycle-integration.txt
  - benchmark-results/m041-s02-t06-go-test.txt
  - benchmark-results/m041-s02-t06-lint.txt
  - benchmark-results/m041-s02-t06-govulncheck.txt
key_decisions:
  - Added `FD_MAX_IN_FLIGHT` as an operator-controlled max in-flight gate for F-2 `model_overloaded`; default 0 preserves unlimited current behavior.
duration: 
verification_result: passed
completed_at: 2026-06-14T06:19:02.604Z
blocker_discovered: false
---

# T06: Added executable lifecycle integration coverage for startup readiness, F-1 model_not_loaded, F-2 model_overloaded, and F-5 shutdown drain behavior.

**Added executable lifecycle integration coverage for startup readiness, F-1 model_not_loaded, F-2 model_overloaded, and F-5 shutdown drain behavior.**

## What Happened

Implemented integration-style lifecycle tests in `api/fd_v2_lifecycle_integration_test.go` because the repository's Go module root is `api/` and the originally planned root-level `tests/integration/...` command is not executable from the current module layout. The tests build the production-style gin chain (`ValidateEmbeddingsRequest` -> `LifecycleGateWithCapacity` -> `EmbeddingsHandler`) with deterministic in-memory embedder/cache fakes. To make F-2 real rather than fabricated, extended lifecycle state with `TryTrackRequest(maxInFlight)`, extended middleware with `LifecycleGateWithCapacity`, and wired `main.go` to operator-controlled `FD_MAX_IN_FLIGHT` (default 0 = unlimited). Covered startup sequence (`/live` 200, `/ready` 503 before warmup, `/ready` 200 after warmup, `/health` still OK), F-1 (`model_not_loaded` + `Retry-After: 5` before warmup then 200), F-2 (`model_overloaded` + `Retry-After: 5` while max-in-flight is exhausted then 200 after capacity recovers), and F-5 (`shutting_down` + `Retry-After: 30` for new requests while an in-flight request drains normally).

## Verification

Fresh lifecycle-specific evidence passed: `go test . -run TestFdV2Lifecycle -v` exit 0. Fresh M043 gate also passed: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. Evidence files: benchmark-results/m041-s02-t06-lifecycle-integration.txt, benchmark-results/m041-s02-t06-go-test.txt, benchmark-results/m041-s02-t06-lint.txt, benchmark-results/m041-s02-t06-govulncheck.txt. GitNexus detect_changes reports MEDIUM risk because `api/main.go` lifecycle capacity config changed, no high/critical risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test . -run TestFdV2Lifecycle -v` | 0 | ✅ pass: startup, F-1, F-2, and F-5 lifecycle integration tests passed | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

The planned path `tests/integration/fd_v2_lifecycle_test.go` was not used because `go test ./tests/integration/...` from the repository root fails: the Go module root is `api/`. The executable equivalent was placed under `api/` and verified with `go test . -run TestFdV2Lifecycle -v`. F-2 required adding a minimal production capacity gate (`FD_MAX_IN_FLIGHT`) because `model_overloaded` previously existed only as a registered error code with no reachable implementation.

## Known Issues

`api/report.json` remains an unrelated untracked generated file outside this task. T06 changes are not committed yet.

## Files Created/Modified

- `api/fd_v2_lifecycle_integration_test.go`
- `api/lifecycle/state.go`
- `api/lifecycle/state_test.go`
- `api/middleware/lifecycle.go`
- `api/middleware/lifecycle_test.go`
- `api/main.go`
- `benchmark-results/m041-s02-t06-lifecycle-integration.txt`
- `benchmark-results/m041-s02-t06-go-test.txt`
- `benchmark-results/m041-s02-t06-lint.txt`
- `benchmark-results/m041-s02-t06-govulncheck.txt`
