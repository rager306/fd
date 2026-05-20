---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Start tagged ONNX 512 service

Start tagged Go ONNX service with `ONNX_MAX_SEQUENCE_LENGTH=512`, isolated Redis namespace, and native HF tokenizer; verify health.

## Inputs

- `api`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Healthy tagged ONNX endpoint on local port`

## Verification

`/health` returns ok for the tagged ONNX service.

## Observability Impact

Runtime health proof for the 512-token path.
