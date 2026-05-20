---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Start ONNX 1024 benchmark service

Start tagged Go ONNX service with `ONNX_MAX_SEQUENCE_LENGTH=1024`, isolated namespace, and native HF tokenizer for benchmarking.

## Inputs

- `api`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Healthy tagged ONNX 1024 benchmark endpoint`

## Verification

`/health` returns ok for the benchmark service.

## Observability Impact

Runtime health proof for benchmark target.
