---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Implement ONNX startup preflight diagnostics

Extend ONNX manifest/runtime config to include validated max sequence length and safe runtime status; fail when configured sequence length exceeds manifest validated contract; log safe metadata.

## Inputs

- `api/main.go`
- `api/embed/onnx_manifest.go`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Updated startup config and tests`

## Verification

Config/manifest tests pass, including sequence length mismatch and safe status fields.

## Observability Impact

Startup preflight reports actionable configuration and manifest diagnostics.
