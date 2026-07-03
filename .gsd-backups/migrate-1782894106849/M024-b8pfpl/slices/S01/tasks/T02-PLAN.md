---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Run packaged ONNX benchmark

Run `benchmark.py` against packaged ONNX Docker image and write `benchmark-results/fd-benchmark-m024-onnx-docker1024.txt` with sanitized effective config.

## Inputs

- `benchmark.py`
- `benchmark-results/fd-environment-inxi-m008.txt`

## Expected Output

- `benchmark-results/fd-benchmark-m024-onnx-docker1024.txt`

## Verification

Benchmark exits 0 and artifact contains packaged runtime config markers.

## Observability Impact

Produces comparable performance artifact with config snapshot.
