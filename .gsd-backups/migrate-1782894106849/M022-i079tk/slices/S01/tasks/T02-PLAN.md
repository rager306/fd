---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Implement ONNX packaging path

Implement the dedicated ONNX Docker packaging path, preferably a root-context Dockerfile plus script that verifies artifacts before build and keeps binaries untracked.

## Inputs

- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`

## Verification

Script shell syntax passes and artifact verifier is invoked by the script.

## Observability Impact

Creates an auditable opt-in packaging path with explicit checks.
