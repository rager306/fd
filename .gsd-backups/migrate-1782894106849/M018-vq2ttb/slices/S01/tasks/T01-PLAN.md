---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Prepare 1024 gate command plan

Confirm runtime prerequisites and exact command for the tagged Go ONNX 1024 legal gate.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `Task summary with command plan`

## Verification

Command plan identifies API URLs, namespace, sequence length, and output artifact.

## Observability Impact

Documents exact runtime/evaluator settings before live run.
