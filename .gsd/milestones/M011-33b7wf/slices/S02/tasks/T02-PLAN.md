---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Implement ONNX manifest validation

Add a small ONNX artifact manifest type and validation function in the Go API. It should parse the tracked manifest schema subset, resolve local artifact path, check file existence, size, SHA256, output name `dense_vecs`, expected dimensions `1024`, and `production_default=false`. Include unit tests for valid manifest, missing file, checksum mismatch, invalid output, and invalid dimensions.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`

## Verification

`cd api && go test ./embed -run 'Test.*ONNX.*|Test.*Manifest.*'` passes.

## Observability Impact

Creates deterministic validation errors before ONNX Runtime load.
