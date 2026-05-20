---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Start tagged ONNX 1024 service

Start tagged Go ONNX service with max sequence length 1024, isolated Redis namespace, and native HF tokenizer; verify health.

## Inputs

- `api`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Healthy tagged ONNX endpoint on local port`

## Verification

`/health` returns ok for the tagged ONNX service.

## Observability Impact

Runtime health proof for 1024-token path.
