---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T02: Implement tagged native tokenizer probe

Implement a minimal build-tagged tokenizer package/probe that imports `github.com/daulet/tokenizers` only under the opt-in tag and can encode fixed probes using the native artifact path. Do not replace runtime ONNX tokenizer yet.

## Inputs

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`

## Expected Output

- `api/embed/hf_tokenizer_native.go`
- `api/embed/hf_tokenizer_native_test.go`

## Verification

Default `go test ./... -short` passes; tagged test compiles/runs when CGO_LDFLAGS points at `.gsd/runtime/tokenizers/linux-amd64`.

## Observability Impact

Provides compile-time proof of native dependency isolation and explicit missing-artifact errors.
