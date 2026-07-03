---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Migrate handler test assertions

Migrate representative handler test assertions to Testify require/assert while preserving response behavior.

## Inputs

- `api/handlers/embeddings_integration_test.go`

## Expected Output

- `api/handlers/embeddings_integration_test.go updated`

## Verification

`cd api && go test ./... -short` passes.

## Observability Impact

Handler test failures become clearer.
