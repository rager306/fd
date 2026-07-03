---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Expose tokenizer metadata from ONNX manifest

Extend manifest validation types to expose tokenizer JSON source file size/sha and add tests for metadata parsing.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `api/embed/onnx_manifest.go`

## Expected Output

- `Updated manifest validation model/tests`

## Verification

Manifest tests pass.

## Observability Impact

Makes tokenizer artifact expectations available to startup preflight.
