---
id: T02
parent: S02
milestone: M046-zqzcu6
key_files:
  - api/middleware/validation.go
  - api/handlers/batch_limits.go
  - api/handlers/batch.go
  - api/handlers/v1batch.go
key_decisions:
  - Use route-level body cap plus handler-specific batch validation instead of forcing batch JSON through `/v1/embeddings` middleware.
duration: 
verification_result: passed
completed_at: 2026-06-14T16:08:28.363Z
blocker_discovered: false
---

# T02: Implemented shared body cap middleware and batch-specific input validation before backend work.

**Implemented shared body cap middleware and batch-specific input validation before backend work.**

## What Happened

Added `middleware.LimitRequestBody()` for JSON endpoints with non-embeddings request shapes. Added `maxBatchInputChars` and `maxLegacyBatchInputs`, validated legacy `/embeddings/batch` inputs and v1 nested batch inputs before any cache or TEI calls, and handled `http.MaxBytesError` as `payload_too_large`. Refactored legacy batch validation into `validateLegacyBatchRequest` to keep cyclomatic complexity within lint limits.

## Verification

`cd api && go test ./middleware ./handlers` passed; full `cd api && go test ./...` passed across 9 packages; lint passed after complexity/name cleanup.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./middleware ./handlers` | 0 | ✅ pass | 12000ms |
| 2 | `cd api && go test ./...` | 0 | ✅ pass | 12100ms |
| 3 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...` | 0 | ✅ pass | 8500ms |

## Deviations

Did not reuse `ValidateEmbeddingsRequest` directly because `/v1/batch` and `/embeddings/batch` have incompatible JSON shapes; instead shared only the body cap and added handler-specific shape validation.

## Known Issues

N+1 TEI call behavior is intentionally left to S03 after inputs are bounded.

## Files Created/Modified

- `api/middleware/validation.go`
- `api/handlers/batch_limits.go`
- `api/handlers/batch.go`
- `api/handlers/v1batch.go`
