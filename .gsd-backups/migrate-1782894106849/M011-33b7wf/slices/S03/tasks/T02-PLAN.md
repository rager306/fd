---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Implement ONNX dense embedder or blocker

Add an ONNX dense embedder package implementation behind the existing `handlers.Embedder` interface. Use `sugarme/tokenizer` to load local tokenizer JSON and `yalue/onnxruntime_go` to load/run `dense_vecs`. The implementation should validate manifest first, initialize ONNX Runtime with an explicit shared-library path, tokenize each input, run CPU EP, return normalized 1024-dimensional `[]float32`, and close native resources. If cgo/runtime blocks implementation, record structured blocker instead of stubbing fake embeddings.

## Inputs

- `api/embed/onnx_manifest.go`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `tei-models/deepvk--USER-bge-m3/tokenizer.json`

## Expected Output

- `api/embed/onnx.go`
- `api/embed/onnx_test.go`
- `api/go.mod`
- `api/go.sum`

## Verification

Focused ONNX embedder tests pass or blocker artifact explains exact dependency/runtime failure.

## Observability Impact

Provides explicit ONNX load/run errors and artifact metadata for future startup logs.
