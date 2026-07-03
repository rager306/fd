---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Implement tagged ONNX tokenizer path

Implement tokenizer abstraction and tagged native implementation so untagged builds keep current behavior and tagged builds use HF native tokenizer for ONNX input IDs/masks.

## Inputs

- `api/embed/onnx.go`
- `api/embed/hf_tokenizer_native.go`

## Expected Output

- `api/embed/onnx.go`
- `api/embed/onnx_tokenizer_default.go`
- `api/embed/onnx_tokenizer_hf.go`

## Verification

Default tests pass; tagged tests pass with `CGO_LDFLAGS`.

## Observability Impact

Tagged ONNX startup should expose tokenizer implementation through tests/loggable metadata if possible.
