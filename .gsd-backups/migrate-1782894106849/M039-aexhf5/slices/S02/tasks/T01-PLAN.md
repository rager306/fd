---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Run packaged legal gate

Start packaged image `fd-api:onnx1024-m039` with `ONNX_RUNTIME_SHA256` and isolated legal namespace, run legal retrieval evaluator against TEI/default and packaged ONNX endpoints, stop container, and verify artifact safety.

## Inputs

- `benchmark-results/fd-onnx-docker-smoke-m039-s01.txt`
- `tools/evaluate_legal_retrieval.py`

## Expected Output

- `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt`

## Verification

Legal evaluator exits 0, metrics pass, artifact contains no raw legal text/secrets/signed URLs, container stopped or prepared for next task.

## Observability Impact

Captures packaged legal parity metrics and sanitized effective config.
