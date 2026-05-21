---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Run Go target-runtime performance gate

Start Go ONNX API with a fresh isolated namespace and run a bounded performance driver against the actual Go ONNX endpoint. Stop server after run.

## Inputs

- `benchmark.py`
- `benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt`

## Expected Output

- `benchmark-results/fd-benchmark-m038-go-onnx-target-runtime.txt`

## Verification

Performance benchmark passes or records blocker; sanitized config present; server stopped.

## Observability Impact

Adds performance evidence through actual Go ONNX endpoint with sanitized config snapshot.
