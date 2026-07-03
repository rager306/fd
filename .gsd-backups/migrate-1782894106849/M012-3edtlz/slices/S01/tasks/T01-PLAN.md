---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Design tokenizer baseline artifact format

Inspect the existing dense comparator probe structure and design a tokenizer-baseline artifact format that records safe metadata, token IDs or token hashes, and raw-text exclusion policy.

## Inputs

- `tools/compare_dense_embeddings.py`
- `benchmark-results/fd-dense-comparator-m010-s02.txt`

## Expected Output

- `Task summary with artifact format decision`

## Verification

Artifact schema is documented in task summary and excludes raw probe text.

## Observability Impact

Defines safe diagnostics before writing generator code.
