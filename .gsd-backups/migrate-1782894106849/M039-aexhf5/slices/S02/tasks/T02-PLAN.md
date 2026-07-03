---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Run packaged performance benchmark

Start packaged image `fd-api:onnx1024-m039` with `ONNX_RUNTIME_SHA256` and isolated benchmark namespace, run benchmark.py against packaged ONNX endpoint, stop container, and verify artifact safety.

## Inputs

- `benchmark.py`

## Expected Output

- `benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt`

## Verification

Benchmark exits 0, artifact metrics present, artifact contains no raw text/secrets/signed URLs, container stopped, port clean.

## Observability Impact

Captures packaged performance metrics and sanitized config.
