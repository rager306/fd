---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Run Go target-runtime legal gate

Start Go ONNX API again with a fresh isolated namespace and run selected Russian/legal retrieval evaluator against TEI API on 8000 and Go ONNX API on 18000. Stop server after run.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `tests/44-FZ-2026-articles.jsonl`
- `benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt`

## Expected Output

- `benchmark-results/fd-legal-retrieval-m038-go-onnx-target-runtime.txt`

## Verification

Legal evaluator passes or records blocker; raw legal text not logged; server stopped.

## Observability Impact

Expands target-runtime proof from smoke to legal retrieval through actual Go endpoints.
