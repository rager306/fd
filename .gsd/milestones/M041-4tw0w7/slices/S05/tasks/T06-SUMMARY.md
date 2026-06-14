---
id: T06
parent: S05
milestone: M041-4tw0w7
key_files:
  - api/observability/traces.go
  - api/observability/traces_test.go
  - api/main.go
  - benchmark-results/m041-s05-t06-go-test.txt
  - benchmark-results/m041-s05-t06-lint.txt
  - benchmark-results/m041-s05-t06-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T08:19:52.999Z
blocker_discovered: false
---

# T06: Added `/v1/traces` in-memory request trace ring buffer and endpoint.

**Added `/v1/traces` in-memory request trace ring buffer and endpoint.**

## What Happened

Implemented `observability.TraceStore`, a concurrency-safe ring buffer retaining the last 100 request traces by default. The trace middleware records timestamp, latency_ms, status, model_id, request_id, path, and dimensions after each request, using request/response headers populated by the existing headers middleware. `FD_TRACES_ENABLED` defaults to true and can disable recording. `GET /v1/traces` returns the current snapshot oldest-first and is skipped by the recorder so reading traces does not mutate the returned list. Wired the trace middleware after headers and before auth/rate/handler execution so request IDs are available and rejected requests can still be observed.

## Verification

Targeted tests passed: after 5 requests `/v1/traces` returns 5 entries with timestamp, latency_ms, status, model_id, request_id, path, and dimensions; ring buffer keeps the last N entries; disabled store records nothing. Fresh full M043 gate passed after the changes: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. GitNexus detect_changes reports LOW risk for tracked changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./observability -run 'TestTrace' -v` | 0 | ✅ pass: trace store and endpoint tests pass | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

None.

## Known Issues

Trace storage is in-memory per process and intentionally small (last 100 requests). `api/report.json` and `.gsd/.../S04-CONTINUE.md` remain unrelated untracked files.

## Files Created/Modified

- `api/observability/traces.go`
- `api/observability/traces_test.go`
- `api/main.go`
- `benchmark-results/m041-s05-t06-go-test.txt`
- `benchmark-results/m041-s05-t06-lint.txt`
- `benchmark-results/m041-s05-t06-govulncheck.txt`
