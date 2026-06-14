---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Added red tests proving batch endpoints still amplify TEI calls per input.

Add failing tests that prove `/v1/batch` and `/embeddings/batch` should call the embedder once per bounded miss group rather than once per input. Include ordering checks and cache hit/miss behavior where cheap.

## Inputs

- `.gsd/milestones/M046-zqzcu6/slices/S02/S02-SUMMARY.md`
- `api/handlers/v1batch.go`
- `api/handlers/batch.go`
- `api/handlers/embeddings.go`

## Expected Output

- `api/handlers/v1batch_test.go`
- `api/handlers/embeddings_integration_test.go`

## Verification

cd api && go test ./handlers -run 'Test(V1BatchUsesSingleEmbedCallPerInnerBatch|TestCreateBatchEmbeddingsUsesSingleEmbedCallForMisses)'

## Observability Impact

Tests create executable evidence for issue #3 P1 #4/#5.
