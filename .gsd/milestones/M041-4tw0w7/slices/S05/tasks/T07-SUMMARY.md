---
id: T07
parent: S05
milestone: M041-4tw0w7
key_files:
  - api/openapi/spec.go
  - api/openapi/spec_test.go
  - api/handlers/openapi.go
  - api/handlers/docs.go
  - api/handlers/openapi_test.go
  - api/main.go
  - benchmark-results/m041-s05-t07-openapi-validator.txt
  - benchmark-results/m041-s05-t07-go-test.txt
  - benchmark-results/m041-s05-t07-lint.txt
  - benchmark-results/m041-s05-t07-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T08:27:47.583Z
blocker_discovered: false
---

# T07: Added programmatic OpenAPI 3.1 JSON schema and Swagger UI docs endpoints.

**Added programmatic OpenAPI 3.1 JSON schema and Swagger UI docs endpoints.**

## What Happened

Implemented `api/openapi.Spec()` as a programmatic OpenAPI 3.1 document covering `/health`, `/live`, `/ready`, `/warmup`, `/version`, `/info`, `/metrics`, `/v1/embeddings`, `/v1/batch`, `/v1/healthcheck`, `/v1/traces`, `/openapi.json`, and `/docs`, including request/response schemas, common response headers, error envelope schema, and bearer auth scheme. Added `NewOpenAPIHandler` for `GET /openapi.json` and `NewDocsHandler` for `GET /docs` with Swagger UI via CDN. Wired both routes into `main.go`; auth skip already covers `/openapi.json` and `/docs`.

## Verification

Targeted tests passed for spec shape, required paths/schemas, `/openapi.json` JSON content, and `/docs` Swagger UI HTML. External validation passed via disposable runner: `uvx --from openapi-spec-validator openapi-spec-validator /tmp/fd-openapi.json` exit 0 with `/tmp/fd-openapi.json: OK`. Fresh full M043 gate passed after the changes: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. GitNexus detect_changes reports LOW risk for tracked changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./openapi ./handlers -run 'TestSpec|TestOpenAPI|TestDocs' -v` | 0 | ✅ pass: OpenAPI/docs tests pass | 180000ms |
| 2 | `cd /root/fd/api && uvx --from openapi-spec-validator openapi-spec-validator /tmp/fd-openapi.json` | 0 | ✅ pass: OpenAPI 3.1 document validates | 300000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 5 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

Used `uvx --from openapi-spec-validator` because the Python module was not preinstalled; this avoided adding project dependencies while still running the requested validator.

## Known Issues

Swagger UI assets are loaded from CDN. `api/report.json` and `.gsd/.../S04-CONTINUE.md` remain unrelated untracked files.

## Files Created/Modified

- `api/openapi/spec.go`
- `api/openapi/spec_test.go`
- `api/handlers/openapi.go`
- `api/handlers/docs.go`
- `api/handlers/openapi_test.go`
- `api/main.go`
- `benchmark-results/m041-s05-t07-openapi-validator.txt`
- `benchmark-results/m041-s05-t07-go-test.txt`
- `benchmark-results/m041-s05-t07-lint.txt`
- `benchmark-results/m041-s05-t07-govulncheck.txt`
