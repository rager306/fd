---
id: T02
parent: S03
milestone: M046-zqzcu6
key_files:
  - api/handlers/batch_backend.go
  - api/handlers/batch.go
  - api/handlers/v1batch.go
  - api/handlers/v1batch_test.go
  - api/handlers/embeddings_integration_test.go
key_decisions:
  - Reuse the proven `/v1/embeddings` cache miss collection pattern for batch endpoints without changing the cache interface.
duration: 
verification_result: passed
completed_at: 2026-06-14T16:22:22.777Z
blocker_discovered: false
---

# T02: Implemented batch cache-miss chunking for `/v1/batch` and `/embeddings/batch`.

**Implemented batch cache-miss chunking for `/v1/batch` and `/embeddings/batch`.**

## What Happened

Added a shared batch backend helper that scans request texts with `GetIfPresent`, collects cache misses, calls the embedder once per bounded miss chunk, backfills the cache with `Set`, and restores response order. `/embeddings/batch` now loads vectors in chunks of 32 and then applies its legacy base64-by-default response encoding. `/v1/batch` now loads each bounded inner batch with one embedder call and preserves nested output order. The tests also prove a repeated identical request is served from cache without additional embedder calls.

## Verification

Focused red tests are now green. `cd api && go test ./handlers` passed and `cd api && go test ./...` passed across all packages.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./handlers -run 'Test(V1BatchUsesSingleEmbedCallPerInnerBatch|CreateBatchEmbeddingsUsesSingleEmbedCallForMisses)'` | 0 | ✅ pass | 10000ms |
| 2 | `cd api && go test ./handlers` | 0 | ✅ pass | 28000ms |
| 3 | `cd api && go test ./...` | 0 | ✅ pass | 12000ms |

## Deviations

The helper is package-local to handlers instead of changing the cache interface; this keeps the change small and avoids coupling S03 to S05 LocalCache redesign.

## Known Issues

Runtime UAT and quality gates remain in T03/T04.

## Files Created/Modified

- `api/handlers/batch_backend.go`
- `api/handlers/batch.go`
- `api/handlers/v1batch.go`
- `api/handlers/v1batch_test.go`
- `api/handlers/embeddings_integration_test.go`
