---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Started and verified the tagged ONNX benchmark server on port 18000.

Start tagged ONNX API on port 18000 with `hf_tokenizers`, isolated Redis namespace, and runtime env. Capture startup duration and memory/RSS where practical.

## Inputs

- `api`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`

## Expected Output

- `Task summary with startup and health evidence`

## Verification

Health endpoint returns ok on port 18000 and process metadata is captured.

## Observability Impact

Captures operational startup and health signals for tagged ONNX runtime.
