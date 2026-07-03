---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify and close milestone

Run final project guardrails, validate and complete M032, checkpoint DB, commit locally, run GitNexus reindex/detect, and report state.

## Inputs

- `tools/verify_onnx_export_contract.py`
- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `benchmark-results/fd-onnx-export-contract-verifier-m032-s02.txt`

## Expected Output

- `Local commit`

## Verification

Verifier positive/negative checks, py_compile, Go checks, actionlint, docs leak checks, tracked binary hygiene, GitNexus detect, commit, reindex.

## Observability Impact

Leaves project clean for next milestone.
