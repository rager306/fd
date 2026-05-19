---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Compare ONNX output to TEI baseline or record blocker

If export succeeds, load the ONNX with ONNX Runtime CPU EP, run the same fixed probes as S02, and compare dimensions, finite values, L2 norms, vector hashes, and cosine similarity against the TEI baseline. If export failed, synthesize the blocker evidence instead. Save a concise tracked benchmark/result artifact under `benchmark-results/`.

## Inputs

- `tools/compare_dense_embeddings.py`
- `benchmark-results/fd-dense-comparator-m010-s02.txt`
- `.gsd/runtime/onnx/m010-s03/`

## Expected Output

- `tools/compare_onnx_dense_embeddings.py`
- `benchmark-results/fd-onnx-fp32-m010-s03.txt`

## Verification

Comparison artifact exists and states PASS/FAIL/BLOCKED with output shape/hash/cosine evidence; raw probe texts are not emitted.

## Observability Impact

Produces ONNX-vs-TEI comparison evidence or a clear export/load blocker for S04.
