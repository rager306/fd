---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Verify workflow alignment

Verify workflow behavior locally through actionlint and by checking provisioning still exposes manifest-derived runtime sha in dry-run.

## Inputs

- `.github/workflows/onnx-packaging.yml`
- `tools/provision_onnx_artifacts.py`

## Expected Output

- `Task summary`

## Verification

actionlint, dry-run expected runtime sha, py_compile pass.

## Observability Impact

Confirms workflow input behavior matches helper contract.
