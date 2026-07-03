---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Resolve worst divergence IDs

Extract worst document/query IDs and hashes from `benchmark-results/fd-legal-retrieval-m015-s03.txt`, resolve them back to `tests/44-FZ-2026-articles.jsonl`, and verify IDs/hashes match without printing raw text.

## Inputs

- `benchmark-results/fd-legal-retrieval-m015-s03.txt`
- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `Task summary with resolved ID evidence`

## Verification

Resolved IDs and text hashes match the M015 artifact.

## Observability Impact

Confirms the diagnostic target set is reproducible and not an artifact parsing error.
