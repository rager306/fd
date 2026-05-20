---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Validate artifact contract

Validate the contract, README, manifests, and tracked binary hygiene.

## Inputs

- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/README.md`

## Expected Output

- `Task summary with validation evidence`

## Verification

Script compile/run, manifest JSON validation, and tracked binary checks pass.

## Observability Impact

Confirms artifact contract is usable and safe to commit.
