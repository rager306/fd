---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Verify target runtime contract

Run S01 verification: manifest JSON, contract marker checks, provisioning/export verifier, actionlint, and GitNexus detect.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`
- `benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt`

## Expected Output

- `Task summary`

## Verification

S01 checks pass and GitNexus scope is low risk.

## Observability Impact

Confirms contract consistency before final closure.
