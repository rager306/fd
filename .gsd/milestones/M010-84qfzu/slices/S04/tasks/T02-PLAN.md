---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Verify M010 spike evidence

Run final evidence checks for M010: Go tests, lint, Python comparator script compile checks, artifact presence checks, raw-probe leakage checks, git status scope, and GitNexus change detection. Prepare milestone validation inputs.

## Inputs

- `tools/compare_dense_embeddings.py`
- `tools/export_user_bge_m3_dense_onnx.py`
- `tools/compare_onnx_dense_embeddings.py`
- `benchmark-results/fd-dense-comparator-m010-s02.txt`
- `benchmark-results/fd-onnx-fp32-m010-s03.txt`

## Expected Output

- `Task summary with verification evidence`

## Verification

Go tests/lint pass; Python scripts compile; artifact parser checks pass; `gitnexus_detect_changes` reports expected scope.

## Observability Impact

Confirms the spike is reproducible and did not introduce runtime regressions before milestone validation.
