---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Run TEI baseline benchmark

Run `benchmark.py` against default TEI API with snapshot v2 metadata and write `benchmark-results/fd-benchmark-m014-tei-baseline.txt`. Use Python 3.13 via uv and preserve existing benchmark safety behavior.

## Inputs

- `benchmark.py`

## Expected Output

- `benchmark-results/fd-benchmark-m014-tei-baseline.txt`

## Verification

Benchmark command exits 0 and artifact includes snapshot_version 2.

## Observability Impact

Captures fresh TEI control metrics for comparison with tagged ONNX.
