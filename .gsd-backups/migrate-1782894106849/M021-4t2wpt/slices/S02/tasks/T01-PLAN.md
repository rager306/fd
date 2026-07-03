---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Validate default Docker boundary

Run the default Docker build to prove the TEI/default path still does not require ONNX/native artifacts.

## Inputs

- `api/Dockerfile`

## Expected Output

- `Task summary with Docker build evidence`

## Verification

`docker build -f api/Dockerfile -t fd-api:m021-default api` exits 0, or records a concrete Docker environment blocker.

## Observability Impact

Proves default container path remains independent of opt-in ONNX artifacts.
