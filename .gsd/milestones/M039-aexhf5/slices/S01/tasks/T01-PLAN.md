---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Check packaged runtime prerequisites

Inspect Docker ONNX build script and Dockerfile contract, then verify local required artifacts/paths before build.

## Inputs

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Task summary`

## Verification

Build script/artifact prerequisite checks pass.

## Observability Impact

Avoids stale image or missing artifact confusion.
