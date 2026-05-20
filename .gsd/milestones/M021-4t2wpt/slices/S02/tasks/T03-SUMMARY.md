---
id: T03
parent: S02
milestone: M021-4t2wpt
key_files:
  - api/embed/onnx.go
  - api/embed/onnx_disabled.go
  - api/embed/onnx_types.go
  - api/embed/onnx_token_types.go
  - api/embed/onnx_tokenizer_default.go
  - api/embed/onnx_tokenizer_hf.go
  - api/embed/onnx_test.go
  - tools/verify_onnx_artifacts.py
  - docs/onnx-artifacts/README.md
key_decisions:
  - Default Go/Docker builds now exclude the real ONNX runtime via the `onnx` build tag and use a default stub.
  - Native tokenizer parity can still be tested with `hf_tokenizers`; real ONNX backend tests/builds require `onnx hf_tokenizers`.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:25:50.391Z
blocker_discovered: false
---

# T03: Validated M021 closure after fixing the default Docker ONNX build boundary.

**Validated M021 closure after fixing the default Docker ONNX build boundary.**

## What Happened

Ran fresh closure verification after the packaging boundary fix. The artifact verifier passes, default Go tests pass, pinned GolangCI-Lint reports 0 issues, native tokenizer tagged tests pass, ONNX+native tagged smoke tests pass, default Docker build passes, no ONNX/native binaries are tracked, no background processes remain, port 18000 is clean, and GitNexus reports low scope with no changed symbols.

## Verification

Fresh verification passed: artifact verifier, default Go tests, lint, tagged tokenizer tests, ONNX+native tagged smoke tests, Docker build, binary hygiene, runtime cleanup, and GitNexus scope.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/verify_onnx_artifacts.py && tools/verify_onnx_artifacts.py` | 0 | ✅ pass — m021_artifact_verifier=pass | 8800ms |
| 2 | `cd api && go test ./... -short` | 0 | ✅ pass — Go test: 74 passed in 4 packages | 32300ms |
| 3 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 32200ms |
| 4 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — Go test: 16 passed in 1 package | 32200ms |
| 5 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — Go test: 2 passed in 1 package | 32100ms |
| 6 | `docker build -f api/Dockerfile -t fd-api:m021-default api` | 0 | ✅ pass — Successfully tagged fd-api:m021-default | 32000ms |
| 7 | `tracked binary, bg_shell, port, GitNexus checks` | 0 | ✅ pass — tracked_native_onnx_binaries=0; no background processes; port_18000_clean; GitNexus low scope | 0ms |

## Deviations

During closure verification, lint initially flagged unused shared ONNX token types in the default build. I split `ONNXEmbedderOptions` from token-only ONNX/HF types, then reran all verification successfully.

## Known Issues

Default Go test count changed from 78 to 74 because ONNX embedder tests are now correctly behind the `onnx` build tag. `hf_tokenizers` test count changed from 20 to 16 for the same reason; targeted `onnx hf_tokenizers` smoke tests pass.

## Files Created/Modified

- `api/embed/onnx.go`
- `api/embed/onnx_disabled.go`
- `api/embed/onnx_types.go`
- `api/embed/onnx_token_types.go`
- `api/embed/onnx_tokenizer_default.go`
- `api/embed/onnx_tokenizer_hf.go`
- `api/embed/onnx_test.go`
- `tools/verify_onnx_artifacts.py`
- `docs/onnx-artifacts/README.md`
