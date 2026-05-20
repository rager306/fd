---
id: T01
parent: S02
milestone: M012-3edtlz
key_files:
  - benchmark-results/fd-tokenizer-baseline-m012-s01.txt
  - api/embed/onnx.go
  - api/embed/onnx_test.go
  - tools/compare_tokenizers.py
key_decisions:
  - Compare Go tokenizer output against the S01 baseline by `label`, `input_ids`, and `attention_mask`.
  - Artifact should record `hf_token_count`, `go_token_count`, IDs/mask hashes, equality booleans, first mismatch index, and short prefix/suffix ID samples for debugging, but never raw text.
  - The current Go path to test is `pretrained.FromFile(tokenizer.json)` plus `EncodeSingle(text, true)`, because that is exactly what `api/embed/onnx.go` uses.
duration: 
verification_result: passed
completed_at: 2026-05-20T02:03:09.859Z
blocker_discovered: false
---

# T01: Defined the exact Go-vs-HF tokenizer parity comparison contract for S02.

**Defined the exact Go-vs-HF tokenizer parity comparison contract for S02.**

## What Happened

Defined the S02 tokenizer comparison contract. The authoritative expected output is the S01 Hugging Face baseline artifact, which contains probe labels, token IDs, masks, hashes, and char counts without raw text. The current Go output must be generated through the same path used in production ONNX prototype code: `pretrained.FromFile(tokenizer.json)` and `EncodeSingle(text, true)`. The comparison artifact should be sanitized and deterministic: no raw probe text, mismatch location and hashes by label, and a non-zero exit when parity fails.

## Verification

Read the HF baseline, ONNX tokenizer implementation, ONNX tests, and S02 plan. The summary records compared fields and artifact safety constraints.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `read benchmark-results/fd-tokenizer-baseline-m012-s01.txt` | 0 | ✅ pass — baseline fields identified | 0ms |
| 2 | `read api/embed/onnx.go and api/embed/onnx_test.go` | 0 | ✅ pass — current Go tokenization path identified | 0ms |
| 3 | `read .gsd/milestones/M012-3edtlz/slices/S02/S02-PLAN.md` | 0 | ✅ pass — acceptance criteria confirmed | 0ms |

## Deviations

None.

## Known Issues

Current Go embedder tokenization uses `sugarme/tokenizer` directly; prior manual probe showed mismatch. S02 T02 should make that mismatch durable and machine-readable.

## Files Created/Modified

- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`
- `api/embed/onnx.go`
- `api/embed/onnx_test.go`
- `tools/compare_tokenizers.py`
