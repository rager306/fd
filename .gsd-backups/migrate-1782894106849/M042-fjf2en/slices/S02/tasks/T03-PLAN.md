---
estimated_steps: 1
estimated_files: 6
skills_used: []
---

# T03: Removed active ONNX build-tagged embedder/runtime files, ONNX Dockerfile, and unused ONNX/tokenizer dependencies from the default module.

Remove ONNX-only build artifacts, Dockerfile paths, module dependencies, and tests that are no longer part of active product scope, or quarantine them in documentation-only research artifacts if deletion is unsafe. Run `go mod tidy` if dependencies are removed. Preserve historical benchmark/docs files.

## Inputs

- `documents/onnx-deactivation-inventory-m042.md`

## Expected Output

- `api/go.mod`
- `api/go.sum`

## Verification

Default `go test ./...` works without ONNX runtime/toolchain dependencies; `go list -deps ./...` no longer includes ONNX runtime packages unless justified.

## Observability Impact

Reduces build/runtime dependency noise.
