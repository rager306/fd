---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T04: /v1/batch endpoint

api/handlers/v1batch.go: POST /v1/batch принимает {"batches": [[s1, s2, ...], [s1, s2, ...], ...]} (каждый inner array ≤ 32 strings, total ≤ 100 batches). Возвращает {"batches": [[e1, e2, ...], [e1, e2, ...], ...]}. Validation аналогично /v1/embeddings. Каждый inner batch обрабатывается через тот же pipeline (validation → lifecycle → cache → model).

## Inputs

- None specified.

## Expected Output

- `api/handlers/v1batch.go`
- `api/handlers/v1batch_test.go`

## Verification

Unit tests: 2 batches × 4 strings → 200 с 2 batches × 4 embeddings. Oversized inner batch → 413. Empty batches → 400 input_required.
