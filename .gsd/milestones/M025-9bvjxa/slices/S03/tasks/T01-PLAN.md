---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Add manual ONNX packaging workflow skeleton

Add `.github/workflows/onnx-packaging.yml` as a manual `workflow_dispatch` skeleton with explicit artifact inputs, provisioning, strict verification, tagged tests, and Docker image build.

## Inputs

- `tools/provision_onnx_artifacts.py`
- `tools/verify_onnx_artifacts.py`
- `tools/build_onnx_image.sh`

## Expected Output

- `.github/workflows/onnx-packaging.yml`

## Verification

Workflow has workflow_dispatch only and actionlint passes.

## Observability Impact

Provides phase-separated CI logs for future ONNX packaging proof.
