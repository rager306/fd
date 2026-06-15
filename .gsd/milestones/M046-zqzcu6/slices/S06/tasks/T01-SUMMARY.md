---
id: T01
parent: S06
milestone: M046-zqzcu6
key_files:
  - api/handlers/embeddings_integration_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T05:43:50.689Z
blocker_discovered: false
---

# T01: Added red test proving `/v1/embeddings` does not yet use batched cache peeks.

**Added red test proving `/v1/embeddings` does not yet use batched cache peeks.**

## What Happened

Added `TestCreateEmbeddingUsesBatchedCachePeek`, which wires a mock cache with `GetManyIfPresent` and asserts the handler performs one batched peek, zero per-item `GetIfPresent` calls, preserves hit/miss ordering, and embeds only cache misses. The test failed as expected because the current handler never calls `GetManyIfPresent`.

## Verification

`cd api && go test ./handlers -run TestCreateEmbeddingUsesBatchedCachePeek` failed with `GetManyIfPresent calls = 0, want 1`, confirming residual P1 #6 before implementation.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./handlers -run TestCreateEmbeddingUsesBatchedCachePeek` | 1 | ✅ expected red fail | 12000ms |

## Deviations

Only the handler-level red test was added in T01. Tiered/Redis batched lookup tests will be added with the implementation in T02 because the public cache API must be extended first.

## Known Issues

Implementation pending in T02.

## Files Created/Modified

- `api/handlers/embeddings_integration_test.go`
