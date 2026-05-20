---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Design opt-in build tag boundary

Design the Go package/build-tag boundary for native HF tokenizers: file names, build tags, interface shape, default fallback behavior, and dependency isolation.

## Inputs

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt`

## Expected Output

- `Task summary with build-tag design`

## Verification

Task summary states exact files/build tags and default-build safety rule.

## Observability Impact

Prevents accidental native dependency in default builds.
