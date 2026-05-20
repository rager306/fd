---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Verify evaluator dry-run hygiene

Run a non-network dry-run/profile mode of the evaluator against the corpus to verify parsing, doc/query derivation, sanitized output, and no raw text leakage.

## Inputs

- `tools/evaluate_legal_retrieval.py`
- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `benchmark-results/fd-legal-retrieval-dry-run-m015-s02.txt`

## Verification

Dry-run artifact exists, includes IDs/counts/hash/thresholds, excludes raw legal text, and GitNexus detect_changes passes.

## Observability Impact

Proves evaluator artifact hygiene before runtime API calls.
