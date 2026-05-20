---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Implement divergence profiler

Create `tools/profile_legal_divergence.py` to compute sanitized token/length/truncation diagnostics for the resolved worst cases using the local HF tokenizer and configurable sequence lengths.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `tools/profile_legal_divergence.py`

## Verification

`python3 -m py_compile tools/profile_legal_divergence.py` passes.

## Observability Impact

Adds repeatable diagnostics for worst-case legal inputs without raw text leakage.
