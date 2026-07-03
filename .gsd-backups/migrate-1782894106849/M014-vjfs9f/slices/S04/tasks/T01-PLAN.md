---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Compare benchmark artifacts

Parse TEI and tagged ONNX benchmark artifacts and write a concise comparison artifact with deltas and caveats.

## Inputs

- `benchmark-results/fd-benchmark-m014-tei-baseline.txt`
- `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt`

## Expected Output

- `benchmark-results/fd-benchmark-m014-comparison.txt`

## Verification

Comparison artifact exists and includes scenario deltas plus caveats.

## Observability Impact

Produces a durable comparison for future runtime decisions.
