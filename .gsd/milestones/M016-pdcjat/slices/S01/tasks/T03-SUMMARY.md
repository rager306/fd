---
id: T03
parent: S01
milestone: M016-pdcjat
key_files:
  - tools/profile_legal_divergence.py
  - benchmark-results/fd-legal-divergence-profile-m016-s01.txt
key_decisions:
  - All 17 worst M015 cases are truncated at 128 tokens by the local HF tokenizer; only 2 remain truncated at 512. This strongly supports max_sequence_length=128 as the primary suspect for S02 runtime diagnostics.
  - S02 should test whether a 512-token ONNX export/runtime aligns with TEI on the same IDs before considering broader remediation.
duration: 
verification_result: passed
completed_at: 2026-05-20T05:35:28.658Z
blocker_discovered: false
---

# T03: Profiled M015 worst cases and found all are truncated at ONNX sequence length 128.

**Profiled M015 worst cases and found all are truncated at ONNX sequence length 128.**

## What Happened

Ran the legal divergence profiler on the M015 worst cases. All 17 worst cases resolved and matched artifact hashes. All 17 are truncated at sequence length 128; only 2 are truncated at 512; maximum tokenized length with specials is 697. The artifact excludes raw legal text and records IDs, hashes, character counts, token counts, sequence-length truncation flags, tokenizer hash, and corpus/M015 artifact hashes.

## Verification

Profiler run passed, hygiene check passed, and GitNexus scope check was low.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece python tools/profile_legal_divergence.py --corpus tests/44-FZ-2026-articles.jsonl --m015-artifact benchmark-results/fd-legal-retrieval-m015-s03.txt --tokenizer-path tei-models/deepvk--USER-bge-m3 --output benchmark-results/fd-legal-divergence-profile-m016-s01.txt --sequence-lengths 128,256,512,1024 --limit 20` | 0 | ✅ pass — 17 cases resolved, all hashes match | 7700ms |
| 2 | `python divergence profile hygiene check` | 0 | ✅ pass — divergence_profile_hygiene=pass; raw_legal_text_leaks=0 | 0ms |
| 3 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ✅ pass — low scope, no changed symbols | 0ms |

## Deviations

None.

## Known Issues

S01 does not prove TEI's internal truncation behavior. It proves the worst cases exceed current ONNX max sequence length 128 under the local HF tokenizer and many fit within 512.

## Files Created/Modified

- `tools/profile_legal_divergence.py`
- `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`
