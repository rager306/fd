---
id: T04
parent: S05
milestone: M041-4tw0w7
key_files:
  - api/handlers/v1batch.go
  - api/handlers/v1batch_test.go
  - api/main.go
  - benchmark-results/m041-s05-t04-go-test.txt
  - benchmark-results/m041-s05-t04-lint.txt
  - benchmark-results/m041-s05-t04-govulncheck.txt
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T08:13:56.468Z
blocker_discovered: false
---

# T04: Added `/v1/batch` endpoint for multiple inner batches with validation, cache/model execution, and tests.

**Added `/v1/batch` endpoint for multiple inner batches with validation, cache/model execution, and tests.**

## What Happened

Implemented `handlers.V1BatchHandler` and wired `POST /v1/batch` in `main.go` behind the lifecycle capacity gate. The endpoint accepts `{"batches":[[...],[...]]}`, requires at least one outer batch, caps outer batches at 100, caps each inner batch at 32 strings, rejects empty inner batches, and returns `{"batches":[[[float...]], ...]}`. Each text is resolved through the existing `EmbeddingCache.GetOrLoad` surface and the configured embedder, preserving fd cache/model behavior. Added tests for 2 batches × 4 strings, oversized inner batch returning 413 `batch_too_large`, and empty batches returning 400 `input_required`.

## Verification

Targeted tests passed for `/v1/batch` success and validation paths. Fresh full M043 gate passed after the changes: `go test ./...` exit 0; golangci-lint v2.12.2 with repo config reports 0 issues; govulncheck exits 0 with 0 reachable vulnerabilities. GitNexus detect_changes reports LOW risk for tracked changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./handlers -run 'TestV1Batch' -v` | 0 | ✅ pass: /v1/batch unit tests pass | 180000ms |
| 2 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go test ./...` | 0 | ✅ pass: all packages ok | 180000ms |
| 3 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass: 0 issues | 180000ms |
| 4 | `cd /root/fd/api && PATH=$PATH:/root/go/bin go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass: 0 reachable vulnerabilities | 300000ms |

## Deviations

The response returns float embeddings for `/v1/batch`; base64 remains the legacy `/embeddings/batch` default and `/v1/embeddings` `encoding_format` option. User-level rate limiting is not applied to `/v1/batch` because this request shape has no `user` field in the task plan.

## Known Issues

`api/report.json` remains an unrelated untracked generated file. `.gsd/.../S04-CONTINUE.md` remains an unrelated auto-compact artifact.

## Files Created/Modified

- `api/handlers/v1batch.go`
- `api/handlers/v1batch_test.go`
- `api/main.go`
- `benchmark-results/m041-s05-t04-go-test.txt`
- `benchmark-results/m041-s05-t04-lint.txt`
- `benchmark-results/m041-s05-t04-govulncheck.txt`
