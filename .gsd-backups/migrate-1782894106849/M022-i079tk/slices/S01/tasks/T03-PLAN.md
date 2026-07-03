---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Run Docker packaging proof

Run the packaging proof: default Docker build, ONNX artifact verification, ONNX image build, and smoke-run `/health` if build succeeds.

## Inputs

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`

## Expected Output

- `Task summary with Docker proof`

## Verification

Default Docker build passes; ONNX build/run passes or records a concrete blocker; cleanup verified.

## Observability Impact

Produces local proof or concrete blocker for ONNX packaging.
