---
id: T03
parent: S02
milestone: M042-fjf2en
key_files:
  - api/go.mod
  - api/go.sum
  - Dockerfile.onnx
  - api/embed/onnx.go
  - api/embed/onnx_disabled.go
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
  - api/embed/onnx_test.go
  - api/embed/onnx_token_types.go
  - api/embed/onnx_tokenizer_default.go
  - api/embed/onnx_tokenizer_hf.go
  - api/embed/onnx_types.go
  - api/embed/hf_tokenizer_native.go
  - api/embed/hf_tokenizer_native_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T10:39:04.968Z
blocker_discovered: false
---

# T03: Removed active ONNX build-tagged embedder/runtime files, ONNX Dockerfile, and unused ONNX/tokenizer dependencies from the default module.

**Removed active ONNX build-tagged embedder/runtime files, ONNX Dockerfile, and unused ONNX/tokenizer dependencies from the default module.**

## What Happened

Deleted the active ONNX Go implementation/test files under `api/embed`, the native HF tokenizer files used by the ONNX path, and `Dockerfile.onnx` after explicit user confirmation. Ran `go mod tidy`, which removed ONNX/tokenizer-related module dependencies from `api/go.mod`/`api/go.sum`. Verified the default module test suite still passes and that `go list -deps ./...` no longer includes ONNX/runtime tokenizer packages.

## Verification

`cd api && go test ./...` passed after deletion and `go mod tidy`. `cd api && go list -deps ./... | rg 'onnx|tokenizers|yalue|daulet'` returned no dependencies.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./...` | 0 | ✅ pass: default module tests pass after ONNX file removal | 180000ms |
| 2 | `cd api && go list -deps ./... | rg 'onnx|tokenizers|yalue|daulet' --color never || true` | 0 | ✅ pass: no ONNX/runtime tokenizer dependencies remain in default dependency graph | 120000ms |

## Deviations

Historical benchmark/GSD/doc artifacts were intentionally preserved. A few ONNX strings remain in health-schema tests/comments and docs; those are not build dependencies and will be addressed in T04 documentation/operator cleanup.

## Known Issues

TEI internal ORT/ONNX probing inside the external TEI binary is not affected by deleting fd's ONNX code. That remains a TEI startup stabilization concern.

## Files Created/Modified

- `api/go.mod`
- `api/go.sum`
- `Dockerfile.onnx`
- `api/embed/onnx.go`
- `api/embed/onnx_disabled.go`
- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`
- `api/embed/onnx_test.go`
- `api/embed/onnx_token_types.go`
- `api/embed/onnx_tokenizer_default.go`
- `api/embed/onnx_tokenizer_hf.go`
- `api/embed/onnx_types.go`
- `api/embed/hf_tokenizer_native.go`
- `api/embed/hf_tokenizer_native_test.go`
