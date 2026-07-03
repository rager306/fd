---
id: T02
parent: S01
milestone: M012-3edtlz
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizer-baseline-m012-s01.txt
key_decisions:
  - Read fixed probes from `tools/compare_dense_embeddings.py` via AST to avoid extra runtime dependencies and keep one source of probe truth.
  - Persist exact `input_ids` and `attention_mask` evidence plus hashes; raw probe text remains excluded.
duration: 
verification_result: passed
completed_at: 2026-05-20T01:59:06.017Z
blocker_discovered: false
---

# T02: Generated the Hugging Face tokenizer baseline artifact for the fixed Russian/legal probes.

**Generated the Hugging Face tokenizer baseline artifact for the fixed Russian/legal probes.**

## What Happened

Added `tools/compare_tokenizers.py`. In baseline mode it loads the local Hugging Face tokenizer from `tei-models/deepvk--USER-bge-m3`, reads the fixed probe list from `tools/compare_dense_embeddings.py` without executing that module, tokenizes each probe with special tokens and attention masks, and writes `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`. The artifact records metadata, package versions, tokenizer/config hashes, probe labels, character counts, token counts, token hashes, and exact token IDs/masks, while excluding raw probe texts.

## Verification

Python compile passed and baseline generation command exited 0 with PASS verdict.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/compare_tokenizers.py && uv run --python 3.13 --with transformers --with torch --with sentencepiece python tools/compare_tokenizers.py --mode baseline` | 0 | ✅ pass — wrote benchmark-results/fd-tokenizer-baseline-m012-s01.txt with PASS verdict | 0ms |

## Deviations

Avoided importing `compare_dense_embeddings.py` directly because it imports `requests`; `compare_tokenizers.py` reads the `PROBES` literal through Python AST instead so tokenizer baseline generation needs only tokenizer dependencies.

## Known Issues

None.

## Files Created/Modified

- `tools/compare_tokenizers.py`
- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`
