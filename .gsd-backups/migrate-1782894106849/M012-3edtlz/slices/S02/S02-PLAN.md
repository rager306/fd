# S02: Go tokenizer candidate comparison

**Goal:** Compare current and candidate Go tokenizer outputs against the S01 Hugging Face baseline and decide whether pure-Go parity is feasible.
**Demo:** After this, current and candidate Go tokenizers are compared against the HF baseline with exact mismatch evidence.

## Must-Haves

- Current `sugarme/tokenizer` output is compared against S01 HF baseline.
- Mismatch artifact records label, lengths, hashes, first mismatch index, and pass/fail.
- At least one alternative/configuration is researched or tested if current tokenizer fails.
- No raw probe text is rendered.
- S03 path is explicit: patch/switch tokenizer, sidecar, or defer.

## Proof Level

- This slice proves: Executable Go/Python comparison evidence and research synthesis.

## Integration Closure

Feeds S03 with an implementation path or a concrete blocker. No ONNX performance work occurs in this slice.

## Verification

- Produces mismatch evidence by label, token length, hash, and first mismatch index without raw probe text.

## Tasks

- [x] **T01: Define Go tokenizer comparison contract** `est:small`
  Inspect S01 baseline structure and current Go ONNX tokenization code to define exact comparison inputs/outputs for Go tokenizer parity.
  - Files: `tools/compare_tokenizers.py`, `api/embed/onnx.go`, `api/embed/onnx_test.go`
  - Verify: Task summary names exact compared fields and artifact safety constraints.

- [x] **T02: Compare current Go tokenizer against HF baseline** `est:medium`
  Add a small Go tokenizer probe/comparison command or test helper that loads the local tokenizer JSON, tokenizes the fixed probes using the current Go path, and emits sanitized JSON for comparison.
  - Files: `tools/compare_tokenizers.py`, `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`
  - Verify: Comparison command exits non-zero if parity fails and writes mismatch artifact without raw text.

- [x] **T03: Evaluate alternative tokenizer paths** `est:medium`
  Research or quickly test alternative pure-Go tokenizer configurations/libraries for XLM-R/SentencePiece parity. Summarize whether S03 should patch current code, switch dependency, use sidecar, or defer.
  - Files: `.gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md`
  - Verify: Research artifact exists and recommends S03 path with evidence.

- [x] **T04: Verify tokenizer comparison artifacts** `est:small`
  Run verification for S02 artifacts: parser checks, raw-text leakage checks, Go/Python compile or tests, and GitNexus detect_changes.
  - Files: `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`, `.gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md`
  - Verify: Artifact parser/leak checks and GitNexus detect_changes pass.

## Files Likely Touched

- tools/compare_tokenizers.py
- api/embed/onnx.go
- api/embed/onnx_test.go
- benchmark-results/fd-tokenizer-go-current-m012-s02.txt
- .gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md
