---
id: T02
parent: S03
milestone: M011-33b7wf
key_files:
  - api/embed/onnx.go
  - api/embed/onnx_test.go
  - api/embed/onnx_manifest.go
  - api/go.mod
  - api/go.sum
key_decisions:
  - Implemented a real Go ONNX embedder with `yalue/onnxruntime_go` and `sugarme/tokenizer`, rather than a Python sidecar or fake stub.
  - ONNX Runtime shared library path is explicit and required; it is not discovered implicitly.
  - Tokenizer path is explicit and required; inference uses `EncodeSingle(text, true)` to include special tokens like Python/HF tokenization.
  - `DynamicAdvancedSession` is protected by a mutex for concurrent handler use until thread-safety is proven.
duration: 
verification_result: passed
completed_at: 2026-05-19T19:18:59.390Z
blocker_discovered: false
---

# T02: Implemented and live-tested a real Go ONNX dense embedder behind explicit artifact, tokenizer, and shared-library inputs.

**Implemented and live-tested a real Go ONNX dense embedder behind explicit artifact, tokenizer, and shared-library inputs.**

## What Happened

Implemented the opt-in ONNX dense embedder in `api/embed/onnx.go`. It validates the manifest, requires an explicit ONNX Runtime shared library path, loads the local tokenizer JSON with `sugarme/tokenizer`, initializes ONNX Runtime through `yalue/onnxruntime_go`, creates a dynamic session for `input_ids` and `attention_mask`, tokenizes/pads batches, runs `dense_vecs`, and returns copied 1024-dimensional float32 vectors. Added tests for constructor error paths and an env-gated live local artifact test. The live test passed using the M010 ONNX artifact, local tokenizer, and a uv-cache `libonnxruntime.so.1.26.0`. Also improved manifest validation so repo-root-relative artifact paths can resolve when tests run from package directories.

## Verification

Focused non-live ONNX/embed tests passed. Env-gated live test passed with local artifact and ONNX Runtime shared library. Dependencies were added to `api/go.mod`/`api/go.sum`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gofmt -w api/embed/onnx_manifest.go api/embed/onnx.go api/embed/onnx_test.go && cd api && go test ./embed -run 'Test.*ONNX.*|Test.*Manifest.*'` | 0 | ✅ pass — ok fd-api/embed | 0ms |
| 2 | `cd api && FD_TEST_ONNX_RUNTIME_LIBRARY=/root/.cache/uv/archive-v0/tj7L7fjW7RE4nnYdVfYZ1/lib/python3.13/site-packages/onnxruntime/capi/libonnxruntime.so.1.26.0 FD_TEST_ONNX_ARTIFACT_MANIFEST=../../docs/onnx-artifacts/user-bge-m3-dense-fp32.json FD_TEST_ONNX_TOKENIZER_PATH=../../tei-models/deepvk--USER-bge-m3/tokenizer.json go test ./embed -run TestONNXEmbedderLiveLocalArtifact -count=1 -v` | 0 | ✅ pass — 1 live ONNX embedder test passed | 6200ms |

## Deviations

Added a manifest path-resolution improvement because tracked manifests use repo-root-relative artifact paths while Go package tests run from package directories. This makes validation robust across working directories.

## Known Issues

Live test used an ONNX Runtime shared library from the uv cache, not a project-managed artifact. Future S04/docs must explain that `ONNX_RUNTIME_LIBRARY` needs a stable path. The embedder uses tokenizer truncation/padding in Go; S03/S04 still need API-level TEI comparison evidence.

## Files Created/Modified

- `api/embed/onnx.go`
- `api/embed/onnx_test.go`
- `api/embed/onnx_manifest.go`
- `api/go.mod`
- `api/go.sum`
