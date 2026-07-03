---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify exact binary contract

Run slice-level verification: manifest JSON validity, docs markers, provisioning dry-run/verifier/export-contract, actionlint, and GitNexus detect.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`
- `benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt`

## Expected Output

- `S01 task summaries`

## Verification

All slice-level checks pass and GitNexus reports expected scope.

## Observability Impact

Confirms the contract is internally consistent without external action.
