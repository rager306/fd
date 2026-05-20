---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Run opt in ONNX API comparison

Run the opt-in ONNX backend against the local artifact if dependencies are viable, compare `/v1/embeddings` output against M010 TEI baseline using existing comparator flow or a small API smoke script, and save tracked evidence. If backend cannot run, save a blocker artifact with exact failure. Verify no raw texts are logged.

## Inputs

- `benchmark-results/fd-dense-comparator-m010-s02.txt`
- `tools/compare_dense_embeddings.py`

## Expected Output

- `benchmark-results/fd-go-onnx-m011-s03.txt`

## Verification

ONNX API comparison artifact exists with PASS/FAIL/BLOCKED and no raw probe text leakage.

## Observability Impact

Creates runtime-level opt-in ONNX evidence or blocker for S04.
