---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Add artifact verification contract

Add a local artifact verification/staging contract that checks ONNX and native tokenizer manifests against local ignored files and emits actionable errors without printing secrets or raw text.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/README.md`

## Verification

`python3 tools/verify_onnx_artifacts.py --onnx-manifest docs/onnx-artifacts/user-bge-m3-dense-fp32.json --native-tokenizer-manifest docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` passes locally.

## Observability Impact

Future Docker/CI work can fail fast on missing or mismatched artifacts.
