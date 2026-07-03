---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify and close milestone

Run final guardrails, complete M034, checkpoint DB, commit locally, reindex GitNexus, and report state.

## Inputs

- `.github/workflows/onnx-packaging.yml`
- `docs/onnx-artifacts/PROVISIONING.md`
- `benchmark-results/fd-onnx-workflow-input-alignment-m034-s02.txt`

## Expected Output

- `Local commit`

## Verification

actionlint, py_compile/provisioning/verifier/export-contract, Go checks, docs leak checks, tracked binary hygiene, GitNexus detect, commit, reindex.

## Observability Impact

Leaves project clean for next milestone.
