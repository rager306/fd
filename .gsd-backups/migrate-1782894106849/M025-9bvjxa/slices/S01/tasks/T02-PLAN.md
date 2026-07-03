---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Implement provisioning helper

Add `tools/provision_onnx_artifacts.py` supporting dry-run and checksum-verified download/copy into manifest local paths. It should require explicit source URLs/paths and never provide fake defaults for missing external artifacts.

## Inputs

- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/*.json`

## Expected Output

- `tools/provision_onnx_artifacts.py`

## Verification

Python compile, dry-run output, missing-source failure behavior, and strict local verifier pass.

## Observability Impact

Gives CI/deploy a reusable pre-build provisioning primitive with actionable failure messages.
