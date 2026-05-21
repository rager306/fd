---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify and close milestone

Run final guardrails, validate and complete M031, checkpoint DB, commit locally, run GitNexus reindex/detect, and report state.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `benchmark-results/fd-onnx-source-contract-m031-s02.txt`

## Expected Output

- `Local commit`

## Verification

Docs/manifests checks, py_compile/provisioning/verifier allow-missing, Go checks as relevant, actionlint, binary hygiene, GitNexus detect, commit, reindex.

## Observability Impact

Leaves project clean for next milestone.
