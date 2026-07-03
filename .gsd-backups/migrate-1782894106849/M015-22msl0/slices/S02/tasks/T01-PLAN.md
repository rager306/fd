---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Implement legal retrieval evaluator

Create `tools/evaluate_legal_retrieval.py` to load the 44-ФЗ JSONL, build sanitized docs/queries, call TEI and ONNX APIs, compute top-k overlap and synthetic known-item metrics, and render a no-raw-text markdown artifact.

## Inputs

- `tests/44-FZ-2026-articles.jsonl`
- `tools/compare_dense_embeddings.py`

## Expected Output

- `tools/evaluate_legal_retrieval.py`

## Verification

`python3 -m py_compile tools/evaluate_legal_retrieval.py` passes.

## Observability Impact

Creates repeatable legal retrieval parity measurement.
