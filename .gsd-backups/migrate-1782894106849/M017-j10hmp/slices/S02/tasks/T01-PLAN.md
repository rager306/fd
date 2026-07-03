---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Write 512 outcome assessment

Write an outcome assessment for the 512-token ONNX gate, comparing M015 128-token failure, M016 Python 512 diagnostic, and M017 tagged Go 512 results.

## Inputs

- `benchmark-results/fd-legal-retrieval-m015-s03.txt`
- `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt`
- `benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt`

## Expected Output

- `benchmark-results/fd-onnx-512-outcome-m017-s02.txt`

## Verification

Artifact exists, includes key metrics, and contains no raw legal text.

## Observability Impact

Durable comparison artifact for future chunking/longer-sequence implementation.
