---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Compare current Go tokenizer against HF baseline

Add a small Go tokenizer probe/comparison command or test helper that loads the local tokenizer JSON, tokenizes the fixed probes using the current Go path, and emits sanitized JSON for comparison.

## Inputs

- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`
- `tei-models/deepvk--USER-bge-m3/tokenizer.json`

## Expected Output

- `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`

## Verification

Comparison command exits non-zero if parity fails and writes mismatch artifact without raw text.

## Observability Impact

Captures current Go tokenizer mismatch evidence in a durable artifact.
