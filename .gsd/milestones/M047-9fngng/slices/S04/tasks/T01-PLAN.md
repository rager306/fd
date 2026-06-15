---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Pinned warmup retry behavior with red tests.

Add tests proving warmup can fail once then retry and mark ready, and bounded terminal failure records last error. Use injectable retry policy with zero delay for deterministic tests.

## Inputs

- `api/main.go`
- `api/lifecycle/warmup.go`
- `documents/issue-6-current-m047.md`

## Expected Output

- `api/main_test.go`

## Verification

cd api && go test ./... (expected red before implementation).

## Observability Impact

Tests pin retry attempt behavior and failure state.
