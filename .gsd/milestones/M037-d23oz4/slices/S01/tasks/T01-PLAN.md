---
estimated_steps: 1
estimated_files: 6
skills_used: []
---

# T01: Inspect target runtime surfaces

Inspect existing Go ONNX runtime tests, API integration tests, legal evaluator, and benchmark harness to ground the target-runtime acceptance contract in actual project surfaces.

## Inputs

- `GitNexus query results`
- `api/embed/onnx_test.go`
- `api/embed/hf_tokenizer_native_test.go`
- `api/main_test.go`
- `tools/evaluate_legal_retrieval.py`
- `benchmark.py`

## Expected Output

- `Task summary`

## Verification

Summarize Go/API/package validation surfaces and Python-helper boundary.

## Observability Impact

Avoids inventing gates not tied to actual runtime surfaces.
