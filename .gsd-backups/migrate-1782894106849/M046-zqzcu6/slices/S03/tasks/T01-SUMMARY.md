---
id: T01
parent: S03
milestone: M046-zqzcu6
key_files:
  - api/handlers/v1batch_test.go
  - api/handlers/embeddings_integration_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T16:19:08.018Z
blocker_discovered: false
---

# T01: Added red tests proving batch endpoints still amplify TEI calls per input.

**Added red tests proving batch endpoints still amplify TEI calls per input.**

## What Happened

Added call-count tests for `/embeddings/batch` and `/v1/batch`. The legacy batch test sends four inputs and expects one embedder call with all four texts. The v1 batch test sends two inner batches of four texts each and expects two embedder calls, one per inner batch. The current implementation failed exactly as expected: legacy made four calls and v1 made eight calls.

## Verification

`cd api && go test ./handlers -run 'Test(V1BatchUsesSingleEmbedCallPerInnerBatch|CreateBatchEmbeddingsUsesSingleEmbedCallForMisses)'` failed with `embedder calls = 4, want 1` and `embedder calls = 8, want 2`, confirming issue #3 P1 #4/#5 behavior before the fix.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./handlers -run 'Test(V1BatchUsesSingleEmbedCallPerInnerBatch|CreateBatchEmbeddingsUsesSingleEmbedCallForMisses)'` | 1 | ✅ expected red fail | 13000ms |

## Deviations

None.

## Known Issues

Implementation still pending in T02.

## Files Created/Modified

- `api/handlers/v1batch_test.go`
- `api/handlers/embeddings_integration_test.go`
