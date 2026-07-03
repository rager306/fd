---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Implemented shared body cap middleware and batch-specific input validation before backend work.

Refactor validation middleware carefully: keep `/v1/embeddings` behavior unchanged, expose or add helpers for request body cap/content-length handling, and add batch-specific validation for legacy `inputs` and v1 nested `batches` shapes.

## Inputs

- `api/middleware/validation.go`
- `api/handlers/batch.go`
- `api/handlers/v1batch.go`

## Expected Output

- `api/middleware/validation.go`

## Verification

Targeted validation tests pass; existing `/v1/embeddings` validation tests still pass.

## Observability Impact

Guardrail failures preserve machine-readable error codes.
