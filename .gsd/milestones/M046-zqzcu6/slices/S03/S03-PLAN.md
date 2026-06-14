# S03: Batch backend work shaping

**Goal:** Reduce batch endpoint backend work amplification by batching cache misses into bounded TEI calls for `/v1/batch` and `/embeddings/batch`, while preserving response shape and S02 guardrails.
**Demo:** A mixed cache-miss batch triggers bounded TEI calls per chunk instead of one TEI call per input.

## Must-Haves

- `/v1/batch` no longer calls TEI once per input on cache misses; it issues one TEI call per bounded inner batch/chunk and preserves nested response order.
- `/embeddings/batch` no longer calls TEI once per input on cache misses; it issues one TEI call per bounded chunk and preserves legacy base64-by-default response shape.
- Cache hits skip TEI, misses are backfilled, and mixed hit/miss batches preserve ordering.
- Existing S02 guardrails and valid batch smokes still pass.
- Go test, lint, govulncheck, and runtime UAT pass.

## Proof Level

- This slice proves: TDD red/green plus runtime executable UAT against rebuilt API.

## Integration Closure

S03 consumes S02 bounded input guardrails and leaves auth posture (S04) and LocalCache correctness (S05) unchanged.

## Verification

- Adds durable evidence in benchmark-results documenting TEI call-count reduction, gates, and runtime smoke results; no new runtime metrics unless an existing log/error surface needs preservation.

## Tasks

- [x] **T01: Added red tests proving batch endpoints still amplify TEI calls per input.** `est:45m`
  Add failing tests that prove `/v1/batch` and `/embeddings/batch` should call the embedder once per bounded miss group rather than once per input. Include ordering checks and cache hit/miss behavior where cheap.
  - Files: `api/handlers/v1batch_test.go`, `api/handlers/embeddings_integration_test.go`
  - Verify: cd api && go test ./handlers -run 'Test(V1BatchUsesSingleEmbedCallPerInnerBatch|TestCreateBatchEmbeddingsUsesSingleEmbedCallForMisses)'

- [x] **T02: Implemented batch cache-miss chunking for `/v1/batch` and `/embeddings/batch`.** `est:1h 30m`
  Refactor batch handlers to use `GetIfPresent`/`Set` miss collection and one `Embed` call per bounded chunk. Preserve legacy `/embeddings/batch` encoding defaults and `/v1/batch` nested response ordering. Handle model wrong-count responses as internal errors.
  - Files: `api/handlers/batch.go`, `api/handlers/v1batch.go`, `api/handlers/embeddings.go`
  - Verify: cd api && go test ./handlers && cd api && go test ./...

- [x] **T03: Verified S03 quality gates and wrote batch backend chunking evidence.** `est:30m`
  Run full Go tests, golangci-lint, govulncheck, and a static proof that batch handlers use `GetIfPresent`/`Set` and no longer use per-input `GetOrLoad` loops. Write S03 evidence artifact.
  - Files: `benchmark-results/m046-s03-batch-backend-chunking.md`, `documents/issue-3-audit-remediation-plan-m046.md`
  - Verify: cd api && go test ./... && cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./... && cd api && go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...

- [x] **T04: Completed runtime UAT and updated R029 for batch backend chunking.** `est:30m`
  Rebuild API container, verify valid `/v1/batch` and `/embeddings/batch` still work, verify S02 rejection paths remain intact, save structured UAT, update requirements/roadmap, and complete S03.
  - Files: `.gsd/milestones/M046-zqzcu6/slices/S03/S03-SUMMARY.md`, `.gsd/milestones/M046-zqzcu6/slices/S03/S03-UAT.md`
  - Verify: docker compose up -d --build api; runtime UAT via gsd_uat_exec; gsd_uat_result_save

## Files Likely Touched

- api/handlers/v1batch_test.go
- api/handlers/embeddings_integration_test.go
- api/handlers/batch.go
- api/handlers/v1batch.go
- api/handlers/embeddings.go
- benchmark-results/m046-s03-batch-backend-chunking.md
- documents/issue-3-audit-remediation-plan-m046.md
- .gsd/milestones/M046-zqzcu6/slices/S03/S03-SUMMARY.md
- .gsd/milestones/M046-zqzcu6/slices/S03/S03-UAT.md
