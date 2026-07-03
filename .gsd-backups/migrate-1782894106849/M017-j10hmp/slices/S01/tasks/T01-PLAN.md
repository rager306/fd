---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Prepare 512 gate command plan

Inspect existing evaluator CLI and current runtime prerequisites for running tagged Go ONNX at max sequence length 512. Confirm required local artifacts exist and identify the exact evaluator command.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `Task summary with command plan`

## Verification

Command plan identifies API URLs, namespace, sequence length, and output artifact.

## Observability Impact

Documents exact runtime/evaluator settings before the live gate.
