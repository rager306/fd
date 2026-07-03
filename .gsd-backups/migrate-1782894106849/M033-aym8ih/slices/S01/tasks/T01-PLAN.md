---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Implement ONNX Runtime wheel extraction

Update `tools/provision_onnx_artifacts.py` to read ONNX Runtime member/size/sha from `source_contract.onnx_runtime`, support safe zip/wheel member extraction, and wire it into `--onnx-runtime-source`. Preserve direct-file fallback and existing tar extraction behavior.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `tools/provision_onnx_artifacts.py`

## Expected Output

- `tools/provision_onnx_artifacts.py`

## Verification

py_compile and targeted synthetic positive probe pass.

## Observability Impact

Adds actionable provisioning failure modes for runtime wheel sources.
