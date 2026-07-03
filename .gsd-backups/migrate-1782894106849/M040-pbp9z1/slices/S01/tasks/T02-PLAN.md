---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Expose TEI runtime metadata in health

Update the runtime health path so TEI/default reports safe runtime metadata from `/health` comparable to ONNX basics: backend, configured model, dimensions, production/default flag, and cache namespace. Preserve ONNX safe metadata behavior and avoid exposing paths, signed URLs, tokens, raw input text, or secrets. Add/update health tests for TEI metadata.

## Inputs

- `api/main.go`
- `api/handlers/health.go`
- `api/handlers/health_test.go`
- `.gsd/milestones/M040-pbp9z1/slices/S01/S01-RESEARCH.md`

## Expected Output

- `api/main.go`
- `api/handlers/health.go`
- `api/handlers/health_test.go`

## Verification

cd api && go test ./... -short

## Observability Impact

Makes `/health` useful for same-host clients to verify the configured backend/model/cache namespace without leaking sensitive paths.
