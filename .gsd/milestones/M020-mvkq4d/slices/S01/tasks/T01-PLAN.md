---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Choose metadata shape

Inspect current ONNX manifest fields and decide whether to update the existing manifest or add a dedicated 1024 runtime contract file.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt`
- `benchmark-results/fd-benchmark-m019-onnx1024.txt`

## Expected Output

- `Task summary with chosen metadata shape`

## Verification

Chosen shape explicitly separates export provenance from validated runtime sequence length.

## Observability Impact

Prevents ambiguity between export-time trace length and validated runtime length.
