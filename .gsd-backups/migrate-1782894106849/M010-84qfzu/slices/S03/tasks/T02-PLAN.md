---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Attempt FP32 dense-only ONNX export

Implement and run a model-preserving FP32 dense-only ONNX export attempt using the local model snapshot. Prefer an explicit wrapper around `AutoModel` that outputs `dense_vecs = normalize(last_hidden_state[:,0])`, with dynamic batch/sequence axes, CPU export, and metadata capture. Store generated ONNX artifacts under `.gsd/runtime/onnx/m010-s03/`.

## Inputs

- `tei-models/deepvk--USER-bge-m3/`

## Expected Output

- `tools/export_user_bge_m3_dense_onnx.py`
- `.gsd/runtime/onnx/m010-s03/export-metadata.json`
- `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`

## Verification

`uv run --python 3.13 --with torch --with transformers --with onnx --with onnxruntime --with safetensors python tools/export_user_bge_m3_dense_onnx.py --model-path tei-models/deepvk--USER-bge-m3 --output-dir .gsd/runtime/onnx/m010-s03` exits 0 or records a structured failure artifact.

## Observability Impact

Captures exact export command, package versions, ONNX file metadata, and failure mode if export cannot complete.
