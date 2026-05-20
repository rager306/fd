---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Design tagged ONNX tokenizer integration

Run impact analysis and design the smallest ONNX tokenizer abstraction needed to swap tokenizers under build tags without changing handlers/cache/API contract.

## Inputs

- `api/embed/onnx.go`
- `api/embed/hf_tokenizer_native.go`

## Expected Output

- `Task summary with integration design`

## Verification

Impact analysis recorded; design names changed symbols and default-build behavior.

## Observability Impact

Keeps runtime integration small and reversible.
