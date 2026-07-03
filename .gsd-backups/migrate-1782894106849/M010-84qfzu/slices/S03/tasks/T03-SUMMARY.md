---
id: T03
parent: S03
milestone: M010-84qfzu
key_files:
  - tools/compare_onnx_dense_embeddings.py
  - benchmark-results/fd-onnx-fp32-m010-s03.txt
  - .gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx
  - .gsd/runtime/onnx/m010-s03/export-metadata.json
key_decisions:
  - Use live TEI/API vectors during ONNX comparison and verify their hashes match the S02 baseline before judging ONNX cosine.
  - Use cosine threshold `0.999` for this FP32 dense-only spike; observed probe cosines were about `0.999993`.
  - Keep ONNX comparator separate from TEI baseline comparator for clear S02/S03 evidence boundaries.
duration: 
verification_result: passed
completed_at: 2026-05-19T18:46:19.461Z
blocker_discovered: false
---

# T03: Compared the exported FP32 dense ONNX candidate against the S02 TEI baseline; all fixed probes passed with cosine around 0.999993.

**Compared the exported FP32 dense ONNX candidate against the S02 TEI baseline; all fixed probes passed with cosine around 0.999993.**

## What Happened

Implemented and ran `tools/compare_onnx_dense_embeddings.py`. The comparator loads the local ONNX artifact with ONNX Runtime CPU EP, tokenizes the same fixed probes with the local tokenizer, requests live TEI/API vectors, verifies live TEI vector hashes match the S02 baseline artifact, and compares TEI vs ONNX dense outputs. The result artifact `benchmark-results/fd-onnx-fp32-m010-s03.txt` passed: all probes had 1024-dimensional finite normalized ONNX output, TEI hashes matched S02, and TEI-vs-ONNX cosine values were all above the `0.999` threshold (observed approximately `0.99999347` to `0.99999393`). Raw probe texts are excluded from the artifact.

## Verification

`py_compile` passed for the ONNX comparator. The comparison command exited 0 and produced a PASS artifact. A parser confirmed required sections/tokens, CPUExecutionProvider/dense_vecs metadata, PASS verdict, raw text logging set false, and no raw probe texts leaked into the artifact.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests --with transformers==4.51.3 --with onnxruntime --with numpy python -m py_compile tools/compare_onnx_dense_embeddings.py` | 0 | ✅ pass | 0ms |
| 2 | `uv run --python 3.13 --with requests --with transformers==4.51.3 --with onnxruntime --with numpy python tools/compare_onnx_dense_embeddings.py --onnx-path .gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx --model-path tei-models/deepvk--USER-bge-m3 --baseline-artifact benchmark-results/fd-dense-comparator-m010-s02.txt --output benchmark-results/fd-onnx-fp32-m010-s03.txt` | 0 | ✅ pass | 6400ms |
| 3 | `python3 artifact validation for benchmark-results/fd-onnx-fp32-m010-s03.txt required sections and raw-probe leakage` | 0 | ✅ pass | 0ms |

## Deviations

None. Export succeeded, so T03 produced a comparison artifact rather than a blocker artifact.

## Known Issues

Comparison uses short fixed probes and verifies dense equivalence, not full retrieval quality. The ONNX artifact is local/ignored and not production-integrated. Export path still relies on `transformers==4.51.3` due failure with latest transformers.

## Files Created/Modified

- `tools/compare_onnx_dense_embeddings.py`
- `benchmark-results/fd-onnx-fp32-m010-s03.txt`
- `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`
- `.gsd/runtime/onnx/m010-s03/export-metadata.json`
