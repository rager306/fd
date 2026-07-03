---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Implement tokenizer runtime provider preflight

Implement startup preflight for tokenizer JSON checksum, optional ONNX Runtime sha, and provider validation; extend health metadata and main tests.

## Inputs

- `api/main.go`
- `api/handlers/health.go`

## Expected Output

- `Updated main config and health metadata`

## Verification

Targeted main/health tests pass.

## Observability Impact

Adds actionable startup diagnostics and safe metadata.
