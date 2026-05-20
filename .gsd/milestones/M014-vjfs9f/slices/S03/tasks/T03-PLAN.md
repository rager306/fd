---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Run tagged ONNX benchmark with restart-aware harness

Run `benchmark.py` against tagged ONNX API with snapshot v3 tagged metadata, configurable API restart command, and write `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`.

## Inputs

- `benchmark.py`

## Expected Output

- `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`

## Verification

Benchmark command exits 0 and artifact includes tagged runtime metadata plus configured restart behavior.

## Observability Impact

Captures ONNX benchmark metrics comparable with S02 TEI baseline.
