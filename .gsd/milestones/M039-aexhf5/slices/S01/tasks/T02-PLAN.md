---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Build packaged ONNX image

Build dedicated ONNX Docker image for M039 from current artifacts and record image id/digest-like metadata.

## Inputs

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Task summary`

## Verification

Docker image build succeeds and image id recorded.

## Observability Impact

Creates packaged artifact for smoke/legal/performance gates.
