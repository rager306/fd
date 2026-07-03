---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Compare HF binding candidate against baseline

Compare the candidate binding output against S01 baseline for all fixed probes and persist a sanitized candidate artifact.

## Inputs

- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`

## Expected Output

- `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt`

## Verification

Candidate comparison exits 0 only if parity passes; otherwise writes mismatch/setup artifact.

## Observability Impact

Proves or disproves candidate token-level parity without raw text.
