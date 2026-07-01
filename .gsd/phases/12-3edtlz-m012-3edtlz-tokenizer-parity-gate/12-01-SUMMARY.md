---
id: S01
parent: M012-3edtlz
milestone: M012-3edtlz
provides:
  - Hugging Face tokenizer baseline artifact.
  - Tokenizer comparison tool foundation for S02.
requires:
  []
affects:
  - S02
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizer-baseline-m012-s01.txt
key_decisions:
  - Use existing `PROBES` from `tools/compare_dense_embeddings.py` as the single fixed probe source.
  - Read `PROBES` via AST to avoid importing API-client dependencies like `requests`.
  - Store exact token IDs/masks and hashes, but never raw probe text.
patterns_established:
  - Use AST extraction to reuse probe literals without importing unrelated dependencies.
  - Separate raw probe text source from safe artifact rendering.
  - Token-level parity should be debugged before embedding cosine differences.
observability_surfaces:
  - `benchmark-results/fd-tokenizer-baseline-m012-s01.txt` provides safe token-level baseline evidence.
  - Task summaries record parser/leak check commands and dependency context.
drill_down_paths:
  - .gsd/milestones/M012-3edtlz/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M012-3edtlz/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M012-3edtlz/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T02:00:43.637Z
blocker_discovered: false
---

# S01: Hugging Face tokenizer baseline

**S01 produced a sanitized Hugging Face tokenization baseline for M012 tokenizer parity work.**

## What Happened

S01 created a deterministic Hugging Face tokenizer baseline for the fixed Russian/legal probes. The new tool loads the local `deepvk/USER-bge-m3` tokenizer, reads the existing comparator probes via AST, tokenizes with Hugging Face behavior, and writes a sanitized artifact containing labels, char counts, token counts, IDs/masks, hashes, package versions, and tokenizer/config hashes. Verification confirmed the artifact parses, records `raw_probe_texts_logged=false`, has five probes, and contains no raw probe strings.

## Verification

Python compile passed; baseline generation exited 0 with PASS verdict; parser/leak check passed; GitNexus detect_changes reported low risk and no affected processes.

## Requirements Advanced

- M011-surfaced-tokenizer-parity — Created the baseline needed to validate tokenizer parity before ONNX performance benchmarking.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Verification helper had to run under `uv` dependency context because `tools/compare_tokenizers.py` imports `transformers`. No artifact or code behavior deviation.

## Known Limitations

The baseline is for the fixed probe set only; larger Russian/legal corpus validation remains future work after fixed-probe parity.

## Follow-ups

S02 should add comparison mode or a Go companion tool to compare current `sugarme/tokenizer` output against this baseline by label, token length, first mismatch index, and hashes.

## Files Created/Modified

- `tools/compare_tokenizers.py` — Tokenizer parity tool with Hugging Face baseline mode.
- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt` — Sanitized HF tokenizer baseline artifact for fixed Russian/legal probes.
