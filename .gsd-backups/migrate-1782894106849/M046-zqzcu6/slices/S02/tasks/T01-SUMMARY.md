---
id: T01
parent: S02
milestone: M046-zqzcu6
key_files:
  - api/handlers/v1batch_test.go
  - api/handlers/embeddings_integration_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T16:08:11.608Z
blocker_discovered: false
---

# T01: Added red tests proving batch endpoints must reject too-long inputs before backend work.

**Added red tests proving batch endpoints must reject too-long inputs before backend work.**

## What Happened

Added tests for `/v1/batch` and `/embeddings/batch` that fail if the embedder is called for an input longer than the batch guardrail limit. The first targeted test run failed before implementation because `maxBatchInputChars` did not exist, establishing a red state.

## Verification

Red test command `cd api && go test ./handlers -run 'Test(V1BatchHandlerRejectsTooLongInputBeforeEmbedder|CreateBatchEmbeddingsRejectsTooLongInputBeforeEmbedder)'` failed before implementation with undefined `maxBatchInputChars`. After implementation, `cd api && go test ./middleware ./handlers` passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./handlers -run 'Test(V1BatchHandlerRejectsTooLongInputBeforeEmbedder|CreateBatchEmbeddingsRejectsTooLongInputBeforeEmbedder)'` | 1 | ✅ expected red fail | 8900ms |
| 2 | `cd api && go test ./middleware ./handlers` | 0 | ✅ pass | 12000ms |

## Deviations

The red test was compile-red rather than behavior-red because the shared limit constant was intentionally introduced by the implementation.

## Known Issues

None for T01.

## Files Created/Modified

- `api/handlers/v1batch_test.go`
- `api/handlers/embeddings_integration_test.go`
