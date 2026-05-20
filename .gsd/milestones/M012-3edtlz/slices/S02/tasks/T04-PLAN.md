---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T04: Verify tokenizer comparison artifacts

Run verification for S02 artifacts: parser checks, raw-text leakage checks, Go/Python compile or tests, and GitNexus detect_changes.

## Inputs

- `tools/compare_tokenizers.py`

## Expected Output

- `Task summary with verification evidence`

## Verification

Artifact parser/leak checks and GitNexus detect_changes pass.

## Observability Impact

Ensures mismatch artifacts are safe and actionable.
