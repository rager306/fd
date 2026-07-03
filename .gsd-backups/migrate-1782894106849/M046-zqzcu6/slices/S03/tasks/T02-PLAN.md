---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Implemented batch cache-miss chunking for `/v1/batch` and `/embeddings/batch`.

Refactor batch handlers to use `GetIfPresent`/`Set` miss collection and one `Embed` call per bounded chunk. Preserve legacy `/embeddings/batch` encoding defaults and `/v1/batch` nested response ordering. Handle model wrong-count responses as internal errors.

## Inputs

- `api/handlers/embeddings.go`
- `api/handlers/batch.go`
- `api/handlers/v1batch.go`

## Expected Output

- `api/handlers/batch.go`
- `api/handlers/v1batch.go`

## Verification

cd api && go test ./handlers && cd api && go test ./...

## Observability Impact

Error logs should include batch/chunk context without logging full input text.
