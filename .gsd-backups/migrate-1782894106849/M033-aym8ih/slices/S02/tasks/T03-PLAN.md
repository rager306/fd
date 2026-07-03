---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify and close milestone

Run final guardrails, complete M033, checkpoint DB, commit locally, run GitNexus reindex/detect, and report state.

## Inputs

- `tools/provision_onnx_artifacts.py`
- `docs/onnx-artifacts/PROVISIONING.md`
- `benchmark-results/fd-onnx-runtime-wheel-provisioning-m033-s02.txt`

## Expected Output

- `Local commit`

## Verification

Synthetic probes, py_compile, Go tests/lint/tagged tests, actionlint, docs leak checks, tracked binary hygiene, GitNexus detect, commit, reindex.

## Observability Impact

Leaves project clean for next milestone.
