---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Profile legal corpus

Profile `tests/44-FZ-2026-articles.jsonl` for schema, counts, invalid flags, length distribution, and corpus hash without printing raw text beyond inspected samples.

## Inputs

- `tests/44-FZ-2026-articles.jsonl`

## Expected Output

- `benchmark-results/fd-legal-corpus-profile-m015-s01.txt`

## Verification

Profile artifact exists and includes counts/hash/no raw text dump.

## Observability Impact

Creates durable corpus metadata for reproducibility.
