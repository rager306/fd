---
id: T01
parent: S03
milestone: M013-nhu1x9
key_files:
  - api/embed/onnx.go
  - api/embed/hf_tokenizer_native.go
key_decisions:
  - Introduce a small `onnxTokenizer` interface used only by `ONNXEmbedder.encodeBatch`.
  - Move default `sugarme/tokenizer` construction behind an untagged helper file, and provide an `hf_tokenizers` tagged helper that returns the parity-correct native tokenizer.
  - Keep `NewONNXEmbedder` signature unchanged so handlers/main/tests remain stable.
  - Risk is MEDIUM for `NewONNXEmbedder` because ONNX tests call it directly; no TEI production path is affected.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:49:34.836Z
blocker_discovered: false
---

# T01: Designed the tagged ONNX tokenizer integration seam with unchanged public constructor and default-build behavior.

**Designed the tagged ONNX tokenizer integration seam with unchanged public constructor and default-build behavior.**

## What Happened

Ran GitNexus impact analysis for `ONNXEmbedder`, `NewONNXEmbedder`, and `encodeBatch`. `ONNXEmbedder` and `encodeBatch` are low risk; `NewONNXEmbedder` is medium risk due direct ONNX tests. The integration design keeps the public constructor and handler contract unchanged. The only intended runtime seam is tokenizer construction and encoding inside the ONNX embedder. Default builds keep current tokenizer behavior; tagged builds use the native HF tokenizer implementation.

## Verification

GitNexus impact analysis completed and design recorded. Risk is limited to ONNX embedder tests and embed module; no TEI production process is affected.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gitnexus_impact target=ONNXEmbedder direction=upstream` | 0 | ✅ reviewed — LOW risk, direct caller NewONNXEmbedder | 0ms |
| 2 | `gitnexus_impact target=NewONNXEmbedder direction=upstream` | 0 | ⚠️ MEDIUM — direct ONNX tests only | 0ms |
| 3 | `gitnexus_impact target=encodeBatch direction=upstream` | 0 | ✅ reviewed — LOW risk, direct caller ONNXEmbedder.Embed | 0ms |

## Deviations

None.

## Known Issues

Integration may still fail if tagged native tokenizer cannot satisfy the exact batch padding/truncation contract used by ONNX; T02/T03 will verify.

## Files Created/Modified

- `api/embed/onnx.go`
- `api/embed/hf_tokenizer_native.go`
