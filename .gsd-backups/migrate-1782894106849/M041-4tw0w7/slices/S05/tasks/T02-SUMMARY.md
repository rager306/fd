---
id: T02
parent: S05
milestone: M041-4tw0w7
key_files:
  - api/middleware/auth.go
  - api/middleware/auth_test.go
  - api/middleware/cors.go
  - api/middleware/cors_test.go
  - api/main.go
  - benchmark-results/m041-s05-t02-go-test.txt
  - benchmark-results/m041-s05-t02-lint.txt
  - benchmark-results/m041-s05-t02-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T08:04:58.586Z
blocker_discovered: false
---

# T02: Added optional FD_API_KEY bearer auth and CORS/preflight middleware with tests and main wiring.

**Added optional FD_API_KEY bearer auth and CORS/preflight middleware with tests and main wiring.**

## What Happened

Implemented `middleware.APIKeyAuthFromEnv`/`APIKeyAuth`: when `FD_API_KEY` is empty auth is disabled; when set, protected endpoints require `Authorization: Bearer <key>` and return OpenAI-style `401 unauthorized` for missing or wrong tokens. Public endpoints are limited to `/live`, `/metrics`, `/docs`, `/docs/*`, and `/openapi.json`; OPTIONS is allowed so CORS preflight can terminate before auth. Implemented `middleware.CORSFromEnv`/`CORS`: default `FD_CORS_ORIGINS` is `*`, comma-separated allowlists are supported, standard methods/headers are emitted, and OPTIONS preflight returns 204. Wired CORS and auth into `main.go` after recovery/headers/metrics so responses remain observable.

## Verification

Targeted tests passed: API key disabled, missing token, wrong token, correct bearer token, public endpoint skip, default CORS preflight 204, allowlisted origin, and disallowed origin behavior. Fresh full M043 gate passed after the changes: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. GitNexus detect_changes reports LOW risk for tracked changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./middleware -run 'TestAPIKeyAuth|TestCORS' -v` | 0 | ✅ pass: auth and CORS unit tests pass | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

None.

## Known Issues

`api/report.json` remains an unrelated untracked generated file. `.gsd/.../S04-CONTINUE.md` remains an unrelated auto-compact artifact.

## Files Created/Modified

- `api/middleware/auth.go`
- `api/middleware/auth_test.go`
- `api/middleware/cors.go`
- `api/middleware/cors_test.go`
- `api/main.go`
- `benchmark-results/m041-s05-t02-go-test.txt`
- `benchmark-results/m041-s05-t02-lint.txt`
- `benchmark-results/m041-s05-t02-govulncheck.txt`
