---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Verify S01 artifact hygiene

Run S01 verification: JSON parser, artifact checksum, default Go tests/lint if code changed, raw/binary leak checks, GitNexus detect_changes.

## Inputs

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `Task summary with verification evidence`

## Verification

All S01 verification gates pass.

## Observability Impact

Confirms artifact contract is safe before build-tag work.
