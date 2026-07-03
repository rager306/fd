---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Added red tests proving batch endpoints must reject too-long inputs before backend work.

Write failing tests that demonstrate `/v1/batch` and `/embeddings/batch` currently accept or reach backend work for oversized bodies, too many inputs, or too-long strings. Include a fake embedder/cache that fails the test if invoked for rejected requests.

## Inputs

- `documents/issue-3-audit-remediation-plan-m046.md`
- `api/main.go`
- `api/handlers/batch.go`
- `api/handlers/v1batch.go`
- `api/middleware/validation.go`

## Expected Output

- `api/handlers/v1batch_test.go`
- `api/handlers/embeddings_integration_test.go`

## Verification

Targeted tests fail before implementation for at least one missing guardrail on each batch endpoint.

## Observability Impact

Tests prove rejected requests do not create backend/cache work.
