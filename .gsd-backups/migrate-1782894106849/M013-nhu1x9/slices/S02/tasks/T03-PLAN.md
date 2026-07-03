---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify tagged tokenizer parity

Run parity test through the tagged package against the M012 baseline and persist a build-tag comparison artifact if possible. If tagged build cannot run, record exact blocker.

## Inputs

- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`

## Expected Output

- `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`

## Verification

Tagged parity command exits 0 and writes PASS artifact, or writes blocker evidence.

## Observability Impact

Proves whether project-local tagged package can produce parity, not just temp-module proof.
