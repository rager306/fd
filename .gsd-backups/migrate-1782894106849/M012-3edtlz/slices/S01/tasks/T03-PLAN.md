---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Verify baseline artifact hygiene

Verify the tokenizer baseline artifact: parse it, check expected sections, ensure raw probe texts are absent, run Python compile, and record evidence.

## Inputs

- `tools/compare_tokenizers.py`
- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`

## Expected Output

- `Task summary with verification evidence`

## Verification

Python compile, artifact parser, raw-text leakage check, and GitNexus detect_changes pass.

## Observability Impact

Guards against unsafe artifact contents and unstable baseline output.
