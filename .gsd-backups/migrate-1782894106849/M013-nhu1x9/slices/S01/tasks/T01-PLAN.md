---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Design native tokenizer artifact contract

Inspect existing artifact-manifest patterns and `.gitignore` coverage for native tokenizer artifacts. Decide the manifest path and local ignored artifact path.

## Inputs

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `.gitignore`

## Expected Output

- `Task summary with manifest path/local path decision`

## Verification

Task summary names tracked manifest path, ignored local artifact path, and binary exclusion rule.

## Observability Impact

Prevents accidental native binary commits and aligns manifest style with ONNX artifact handling.
