---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Wire opt in ONNX backend

Wire the opt-in ONNX backend into `api/main.go`: default TEI path unchanged; when `EMBEDDING_BACKEND=onnx`, validate manifest, require shared library/tokenizer paths, construct ONNX embedder, and route handlers to it. Add tests for default TEI behavior and ONNX config error paths. Do not expose ONNX unless explicitly configured.

## Inputs

- `api/embed/onnx.go`
- `api/embed/onnx_manifest.go`

## Expected Output

- `api/main.go`
- `api/main_test.go`

## Verification

`cd api && go test ./... -short` passes; default env still selects TEI.

## Observability Impact

Startup logs selected backend and ONNX artifact_id/provider on successful opt-in initialization.
