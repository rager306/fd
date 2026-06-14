---
id: T01
parent: S05
milestone: M041-4tw0w7
key_files:
  - api/embed/types.go
  - api/embed/types_test.go
  - api/handlers/embeddings.go
  - api/handlers/embeddings_integration_test.go
  - api/handlers/errors.go
  - api/handlers/errors_test.go
  - api/middleware/validation.go
  - api/middleware/validation_test.go
  - benchmark-results/m041-s05-t01-go-test.txt
  - benchmark-results/m041-s05-t01-lint.txt
  - benchmark-results/m041-s05-t01-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T07:59:06.768Z
blocker_discovered: false
---

# T01: Added OpenAI-compatible `user` and `priority` request fields, priority validation, and tests for base64/user/priority behavior.

**Added OpenAI-compatible `user` and `priority` request fields, priority validation, and tests for base64/user/priority behavior.**

## What Happened

Extended `embed.EmbeddingsRequest` to preserve optional `user` and `priority` fields during custom JSON unmarshalling while keeping existing single-string and array `input` support. Added priority enum validation (`low|normal|high`) in both production middleware and inline handler validation. Added canonical `priority_invalid` error code and registry/test coverage. Existing base64 response support was already present, so T01 adds tests to lock the contract: `/v1/embeddings` with `encoding_format=base64` returns a base64 string, `priority=high` is accepted, `user` is accepted, and invalid priority returns a clean 400 error envelope.

## Verification

Targeted tests passed for embed/middleware/handlers. Fresh full M043 gate passed after the changes: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. GitNexus detect_changes reports LOW risk.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./embed ./middleware ./handlers -run 'TestEmbeddingsRequest|TestValidation|TestCreateEmbedding_ProductionHandler|TestError' -v` | 0 | ✅ pass: targeted base64/user/priority/error tests pass | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

The task expected `api/embed/base64.go`, but base64 encoding already exists in `api/embed/codec.go`; this task validates and uses the existing implementation rather than creating a duplicate file. `api/handlers/batch.go` did not need changes because legacy `/embeddings/batch` already has its own encoding_format support and T01 focuses on `/v1/embeddings` request fields.

## Known Issues

`api/report.json` remains an unrelated untracked generated file. `.gsd/.../S04-CONTINUE.md` remains an unrelated auto-compact artifact.

## Files Created/Modified

- `api/embed/types.go`
- `api/embed/types_test.go`
- `api/handlers/embeddings.go`
- `api/handlers/embeddings_integration_test.go`
- `api/handlers/errors.go`
- `api/handlers/errors_test.go`
- `api/middleware/validation.go`
- `api/middleware/validation_test.go`
- `benchmark-results/m041-s05-t01-go-test.txt`
- `benchmark-results/m041-s05-t01-lint.txt`
- `benchmark-results/m041-s05-t01-govulncheck.txt`
