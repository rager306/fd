---
id: T02
parent: S03
milestone: M010-84qfzu
key_files:
  - tools/export_user_bge_m3_dense_onnx.py
  - .gsd/runtime/onnx/m010-s03/export-metadata.json
  - .gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx
key_decisions:
  - Pin `transformers==4.51.3` for the S03 export path; latest `transformers 5.8.1` is not currently compatible with this legacy torch.onnx trace path.
  - Use explicit dense wrapper `normalize(last_hidden_state[:,0])` instead of exporting a generic SentenceTransformer wrapper.
  - Keep generated ONNX artifact in ignored `.gsd/runtime/onnx/m010-s03/` and track only scripts/results.
duration: 
verification_result: mixed
completed_at: 2026-05-19T18:43:42.278Z
blocker_discovered: false
---

# T02: Exported a local model-preserving FP32 dense-only ONNX candidate for `deepvk/USER-bge-m3` after pinning transformers to 4.51.3.

**Exported a local model-preserving FP32 dense-only ONNX candidate for `deepvk/USER-bge-m3` after pinning transformers to 4.51.3.**

## What Happened

Implemented `tools/export_user_bge_m3_dense_onnx.py`, a local spike tool that loads the exact local `tei-models/deepvk--USER-bge-m3` snapshot, wraps `AutoModel` with dense-only CLS pooling and L2 normalization, exports `dense_vecs` with dynamic batch/sequence axes, and writes structured metadata. The first unpinned dependency run used `transformers 5.8.1` and failed during `torch.onnx.export` tracing with `IndexError: tuple index out of range`, despite dummy forward producing shape `[1, 1024]` and norm ~1.0. A second run pinned `transformers==4.51.3` and succeeded.

Successful export evidence: `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`, size `1432482908` bytes, SHA256 `28538a17a99302e144149732d73fb273cd7c7a0468dc59167caa5a2d5ff2a3d4`; ONNX Runtime CPU EP load succeeded; ONNX inputs are `input_ids` and `attention_mask` with dynamic batch/sequence axes; output is `dense_vecs` with shape `[batch_size, 1024]` and type `tensor(float)`; packages were torch `2.12.0`, transformers `4.51.3`, onnx `1.21.0`, onnxruntime `1.26.0`, safetensors `0.7.0`; production runtime was not changed.

## Verification

`py_compile` passed. Unpinned export failed with a recorded dependency/tracing failure. Pinned export with `transformers==4.51.3` exited 0, wrote `export-metadata.json`, produced a 1.43GB ONNX artifact under ignored runtime storage, and validated ONNX Runtime CPU EP dummy load with output shape `[1,1024]`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with torch --with transformers --with onnx --with onnxruntime --with safetensors python -m py_compile tools/export_user_bge_m3_dense_onnx.py` | 0 | ✅ pass | 0ms |
| 2 | `uv run --python 3.13 --with torch --with transformers --with onnx --with onnxruntime --with safetensors python tools/export_user_bge_m3_dense_onnx.py --model-path tei-models/deepvk--USER-bge-m3 --output-dir .gsd/runtime/onnx/m010-s03` | 1 | ⚠️ failed as dependency/tracing evidence (`transformers 5.8.1`, IndexError) | 16200ms |
| 3 | `uv run --python 3.13 --with torch --with transformers==4.51.3 --with onnx --with onnxruntime --with safetensors python tools/export_user_bge_m3_dense_onnx.py --model-path tei-models/deepvk--USER-bge-m3 --output-dir .gsd/runtime/onnx/m010-s03` | 0 | ✅ pass | 45300ms |
| 4 | `read .gsd/runtime/onnx/m010-s03/export-metadata.json` | 0 | ✅ pass | 0ms |

## Deviations

First export attempt with latest unpinned `transformers 5.8.1` failed with `IndexError: tuple index out of range` during legacy torch.onnx tracing. Retried with `transformers==4.51.3` based on BGE-M3 ONNX reference compatibility; export succeeded. The successful metadata overwrote the failed metadata, but the failure mode is recorded in this task summary.

## Known Issues

The ONNX artifact is local and ignored, not tracked. The export uses legacy TorchScript-based ONNX export and emits a deprecation warning. It has only been dummy-load validated so far; S03 T03 must run real probe comparison against the S02 TEI baseline.

## Files Created/Modified

- `tools/export_user_bge_m3_dense_onnx.py`
- `.gsd/runtime/onnx/m010-s03/export-metadata.json`
- `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`
