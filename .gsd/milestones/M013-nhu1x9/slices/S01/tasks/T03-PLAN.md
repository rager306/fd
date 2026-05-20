---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify native tokenizer artifact contract

Add a small verifier mode or script for native tokenizer artifact manifest validation, and run it against the local temp artifact without committing the binary.

## Inputs

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `Artifact verification evidence`

## Verification

Verifier exits 0 for local artifact; tracked-file scan confirms no `.a` binary is tracked.

## Observability Impact

Makes missing/wrong native artifact failures explicit and repeatable.
