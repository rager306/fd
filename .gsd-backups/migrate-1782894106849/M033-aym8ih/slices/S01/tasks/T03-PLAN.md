---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify provisioning compatibility

Run compatibility checks for existing dry-run and verifier behavior, then record S01 summary.

## Inputs

- `tools/provision_onnx_artifacts.py`
- `tools/verify_onnx_artifacts.py`

## Expected Output

- `Task summary`

## Verification

dry-run/verifier allow-missing and compile checks pass.

## Observability Impact

Ensures existing provisioning contract remains compatible.
