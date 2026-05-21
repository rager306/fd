---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Align workflow runtime sha handling

Update `.github/workflows/onnx-packaging.yml` input descriptions, validation step, and provisioning argument assembly so runtime sha is optional when manifest metadata supplies it, while preserving checksum verification by provisioning.

## Inputs

- `.github/workflows/onnx-packaging.yml`
- `tools/provision_onnx_artifacts.py`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `.github/workflows/onnx-packaging.yml`

## Verification

actionlint and text checks pass.

## Observability Impact

Workflow messages explain checksum source behavior.
