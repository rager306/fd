---
id: T02
parent: S03
milestone: M013-nhu1x9
key_files:
  - api/embed/onnx.go
  - api/embed/onnx_tokenizer_default.go
  - api/embed/onnx_tokenizer_hf.go
  - api/embed/hf_tokenizer_native.go
key_decisions:
  - `ONNXEmbedder` now depends on a small `onnxTokenizer` interface.
  - Default untagged builds use `sugarmeONNXTokenizer` from `onnx_tokenizer_default.go`.
  - Tagged `hf_tokenizers` builds use `newNativeHFTokenizer` through `onnx_tokenizer_hf.go`.
  - `ONNXEmbedder.Close` now closes tokenizer resources as well as ONNX session resources.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:52:55.187Z
blocker_discovered: false
---

# T02: Integrated the tokenizer abstraction so tagged ONNX builds use the HF native tokenizer while default builds keep existing behavior.

**Integrated the tokenizer abstraction so tagged ONNX builds use the HF native tokenizer while default builds keep existing behavior.**

## What Happened

Implemented the tagged ONNX tokenizer path. `ONNXEmbedder` now uses an internal `onnxTokenizer` interface returning input IDs and attention masks. Untagged builds construct a `sugarmeONNXTokenizer`, preserving default behavior. Tagged `hf_tokenizers` builds construct the parity-correct native HF tokenizer. Default tests pass without native flags, and tagged tokenizer tests pass with the validated native library path.

## Verification

Fresh default Go tests passed. Tagged embed tests passed with native tokenizer linker flags.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./... -short` | 0 | ✅ pass — default build preserved | 0ms |
| 2 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -run 'TestNativeHFTokenizerMatchesBaseline|TestONNXEmbedderLiveLocalArtifact' -count=1` | 0 | ✅ pass — tagged native tokenizer test passed; live ONNX test skipped without env vars | 0ms |

## Deviations

Added `api/embed/onnx_tokenizer_default.go` and `api/embed/onnx_tokenizer_hf.go` as planned. The tagged live ONNX artifact test was included in the test regex but skipped because live ONNX env vars were not supplied; T03 will run the real tagged API/cosine path.

## Known Issues

S03 T03 still needs to run a tagged ONNX API instance with real ONNX env vars and verify cosine. Until then, runtime semantic equivalence is not proven.

## Files Created/Modified

- `api/embed/onnx.go`
- `api/embed/onnx_tokenizer_default.go`
- `api/embed/onnx_tokenizer_hf.go`
- `api/embed/hf_tokenizer_native.go`
