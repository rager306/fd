---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Implement TEI API dense comparator

Implement a Python comparator script that queries the current API for the fixed probes, validates output dimensions and finite floats, checks L2 normalization tolerance, computes deterministic vector hashes, and computes pairwise cosine similarity. Keep runtime configurable through env vars like `COMPARE_API_URL`, `COMPARE_MODEL`, and `COMPARE_DIMENSIONS`, with safe defaults matching current fd behavior.

## Inputs

- `benchmark.py`
- `api/handlers/embeddings.go`
- `api/handlers/batch.go`

## Expected Output

- `tools/compare_dense_embeddings.py`

## Verification

`uv run --python 3.13 --with requests python tools/compare_dense_embeddings.py --output benchmark-results/fd-dense-comparator-m010-s02.txt` exits 0 against healthy local stack.

## Observability Impact

Adds reusable local diagnostic tooling for dense embedding equivalence checks.
