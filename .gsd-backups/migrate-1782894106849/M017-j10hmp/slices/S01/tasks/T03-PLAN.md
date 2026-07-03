---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Run full legal quality gate

Run the legal retrieval evaluator against TEI and tagged ONNX 512, then check artifact hygiene and summarize metrics.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt`

## Verification

Evaluator exits 0 or records explicit fail verdict; artifact exists and raw legal text leak check passes.

## Observability Impact

Full legal corpus quality evidence for ONNX 512.
