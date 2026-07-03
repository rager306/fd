---
id: T01
parent: S02
milestone: M021-4t2wpt
key_files:
  - api/embed/onnx.go
  - api/embed/onnx_disabled.go
  - api/embed/onnx_types.go
  - api/embed/onnx_tokenizer_default.go
  - api/embed/onnx_tokenizer_hf.go
  - api/embed/onnx_test.go
  - api/Dockerfile
key_decisions:
  - Real ONNX backend now requires build tag `onnx`; parity-correct native tokenizer still uses `hf_tokenizers`.
  - Default Docker/CI path remains CGO-disabled and does not import `onnxruntime_go`.
  - Validated ONNX runtime commands should use `-tags 'onnx hf_tokenizers'` when using native HF tokenizer.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:22:33.637Z
blocker_discovered: false
---

# T01: Fixed and validated the default Docker build boundary so ONNX remains opt-in.

**Fixed and validated the default Docker build boundary so ONNX remains opt-in.**

## What Happened

Validated the default Docker boundary. The first build exposed a real default-build regression: `onnxruntime_go` was imported without an opt-in ONNX build tag, causing CGO-disabled Docker build failure. I split the ONNX implementation behind `//go:build onnx`, added default stub types, kept shared ONNX option/token types buildable, updated tokenizer build tags, and reran the default Docker build successfully.

## Verification

`docker build -f api/Dockerfile -t fd-api:m021-default api` passed after the build tag fix. Default Go tests and tagged checks were also smoke-run earlier after the fix.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker build -f api/Dockerfile -t fd-api:m021-default api` | 0 | ✅ pass — Successfully tagged fd-api:m021-default | 42100ms |
| 2 | `cd api && go test ./... -short` | 0 | ✅ pass — Go test: 74 passed in 4 packages after ONNX tests moved behind onnx tag | 7200ms |
| 3 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — Go test: 16 passed in 1 package | 7200ms |
| 4 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — Go test: 2 passed in 1 package | 7100ms |

## Deviations

Initial default Docker build failed because `api/embed/onnx.go` imported `onnxruntime_go` in the default CGO-disabled build. I fixed the boundary by putting the real ONNX embedder behind build tag `onnx`, adding a default stub, and preserving native tokenizer tests under `hf_tokenizers`. The retried default Docker build passed.

## Known Issues

Future scripts/docs that start ONNX runtime must use `-tags 'onnx hf_tokenizers'`; M021 README and ONNX manifest runtime env were updated accordingly.

## Files Created/Modified

- `api/embed/onnx.go`
- `api/embed/onnx_disabled.go`
- `api/embed/onnx_types.go`
- `api/embed/onnx_tokenizer_default.go`
- `api/embed/onnx_tokenizer_hf.go`
- `api/embed/onnx_test.go`
- `api/Dockerfile`
