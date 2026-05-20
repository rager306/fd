---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Define benchmark matrix and metadata contract

Inspect current benchmark.py config snapshot and M013 artifacts to decide the minimal metadata additions needed for tagged ONNX benchmark comparability.

## Inputs

- `benchmark.py`
- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Task summary with benchmark matrix and metadata fields`

## Verification

Task summary lists scenarios and metadata fields.

## Observability Impact

Defines the benchmark evidence contract before code changes.
