---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Add tagged ONNX metadata snapshot fields

Extend benchmark.py snapshot to optionally include tagged ONNX/native metadata from env vars and manifest files, while preserving existing TEI behavior and raw text exclusion.

## Inputs

- `benchmark.py`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `benchmark.py`

## Verification

`uv run --python 3.13 --with requests --with redis python -m py_compile benchmark.py` and a lightweight snapshot function check pass.

## Observability Impact

Future benchmark artifacts can record build/runtime mode, native artifact checksum, ONNX artifact checksum, ORT library path/hash, and Redis namespace.
