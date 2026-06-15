---
estimated_steps: 1
estimated_files: 8
skills_used: []
---

# T02: Simplified runtime health, embed interface, and lifecycle state contracts.

Remove ONNX-only RuntimeHealth fields and stale ONNX health tests/fields, define one shared embed interface in `api/embed`, update handlers/lifecycle to use it, and remove lifecycle default singleton in favor of explicit `NewState` in main.

## Inputs

- `api/main.go`
- `api/handlers/health.go`
- `api/lifecycle/state.go`

## Expected Output

- `api/embed/types.go`
- `api/handlers/embeddings.go`
- `api/lifecycle/warmup.go`
- `api/lifecycle/state.go`
- `api/main.go`

## Verification

cd api && go test ./handlers ./lifecycle ./...

## Observability Impact

Health metadata reflects only active TEI runtime fields.
