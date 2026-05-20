---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Run packaged legal retrieval evaluator

Run `tools/evaluate_legal_retrieval.py` against TEI at 8000 and packaged ONNX at 18000, writing a M023 artifact under `benchmark-results/` with runtime labels and cache namespaces.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt`

## Verification

Evaluator exits 0 for pass or nonzero with blocked/fail artifact; artifact excludes raw legal text.

## Observability Impact

Produces sanitized legal gate metrics for packaged runtime.
