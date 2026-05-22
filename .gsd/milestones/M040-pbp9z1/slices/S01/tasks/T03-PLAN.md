---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T03: Resolve embeddings model-field contract

Choose the smallest safe behavior for `/v1/embeddings` request `model`: either harden the handler to reject non-empty model values that do not match the configured model with a 400 and tests, or document that request `model` is compatibility metadata and response model plus `/health.runtime.model` are authoritative. Prefer implementation hardening only if local tests and existing API expectations support it without broad breakage.

## Inputs

- `api/handlers/embeddings.go`
- `api/handlers/embeddings_integration_test.go`
- `docs/same-host-embedding-service-contract.md`
- `.gsd/milestones/M040-pbp9z1/slices/S01/S01-RESEARCH.md`

## Expected Output

- `api/handlers/embeddings.go`
- `api/handlers/embeddings_integration_test.go`
- `docs/same-host-embedding-service-contract.md`

## Verification

cd api && go test ./... -short

## Observability Impact

Prevents ambiguity between client-requested model metadata and the actual configured model served by fd.
