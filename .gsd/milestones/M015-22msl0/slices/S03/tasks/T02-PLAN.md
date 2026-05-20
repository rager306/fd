---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Run live legal retrieval gate

Run `tools/evaluate_legal_retrieval.py` live against TEI and tagged ONNX using the 44-ФЗ corpus and write `benchmark-results/fd-legal-retrieval-m015-s03.txt`.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `benchmark-results/fd-legal-retrieval-m015-s03.txt`

## Verification

Evaluator exits 0 for pass or 2 for quality fail and artifact records verdict.

## Observability Impact

Captures legal retrieval parity metrics and verdict.
