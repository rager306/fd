---
id: T03
parent: S03
milestone: M011-33b7wf
key_files:
  - api/main.go
  - api/main_test.go
  - api/embed/onnx.go
key_decisions:
  - `EMBEDDING_BACKEND=onnx` now requires `ONNX_ARTIFACT_MANIFEST`, `ONNX_RUNTIME_LIBRARY`, and `ONNX_TOKENIZER_PATH`.
  - `ONNX_MAX_SEQUENCE_LENGTH` defaults to 512 and is passed to the embedder.
  - Handlers continue to receive the same `handlers.Embedder` interface, so API response wiring remains unchanged.
duration: 
verification_result: passed
completed_at: 2026-05-19T19:21:31.957Z
blocker_discovered: false
---

# T03: Wired the real ONNX embedder behind explicit opt-in config while preserving TEI as default.

**Wired the real ONNX embedder behind explicit opt-in config while preserving TEI as default.**

## What Happened

Wired the opt-in ONNX backend into `api/main.go`. Default config still creates the TEI client. When `EMBEDDING_BACKEND=onnx`, startup now validates the manifest, requires explicit ONNX Runtime shared library and tokenizer paths, constructs `embed.NewONNXEmbedder`, logs artifact metadata, and passes the ONNX embedder into the existing embedding and batch handlers. Tests were updated for required ONNX env vars, valid config parsing, invalid config handling, and max sequence length parsing. Full Go tests passed.

## Verification

Fresh verification passed: formatted changed Go files and ran `cd api && go test ./... -short`; all packages passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gofmt -w api/main.go api/main_test.go api/embed/onnx.go api/embed/onnx_test.go api/embed/onnx_manifest.go api/embed/onnx_manifest_test.go && cd api && go test ./... -short` | 0 | ✅ pass — api, cache, embed, handlers all ok | 0ms |

## Deviations

S03 replaced the temporary not-implemented ONNX branch from S02 with real ONNX embedder construction. The live API comparison is deferred to T04 as planned.

## Known Issues

The ONNX backend still depends on a stable external shared library path. Current live tests use uv-cache `libonnxruntime.so.1.26.0`; docs/S04 should avoid presenting that as production-ready storage.

## Files Created/Modified

- `api/main.go`
- `api/main_test.go`
- `api/embed/onnx.go`
