---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Run ONNX 1024 benchmark

Run `benchmark.py` against ONNX 1024 with uv/Python 3.13 and save the artifact.

## Inputs

- `benchmark.py`

## Expected Output

- `benchmark-results/fd-benchmark-m019-onnx1024.txt`

## Verification

Benchmark exits 0 and artifact hygiene check passes.

## Observability Impact

Measured ONNX 1024 latency/throughput/cache behavior.
