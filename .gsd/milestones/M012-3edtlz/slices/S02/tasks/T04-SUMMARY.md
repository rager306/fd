---
id: T04
parent: S02
milestone: M012-3edtlz
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizer-go-current-m012-s02.txt
  - .gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md
key_decisions:
  - Treat `go-current` exit code 2 as expected evidence for current mismatch, not verification failure.
  - S02 artifacts are safe: parser/leak check passed and no raw probe text appears in baseline, current-Go comparison, or S02 research.
duration: 
verification_result: passed
completed_at: 2026-05-20T02:11:25.753Z
blocker_discovered: false
---

# T04: Verified S02 tokenizer comparison artifacts and confirmed mismatch evidence is safe and actionable.

**Verified S02 tokenizer comparison artifacts and confirmed mismatch evidence is safe and actionable.**

## What Happened

Ran final S02 verification. The tokenizer comparison tool compiles. Running `go-current` mode exits 2 as expected because the current Go tokenizer does not match the HF baseline, and it writes the mismatch artifact. A parser/leak check verified required sections, `raw_probe_texts_logged=false`, five probes, and zero raw probe text leaks across baseline, comparison, and research artifacts. GitNexus change detection reports medium risk only for internal `tools/compare_tokenizers.py` flows, with no API runtime affected processes.

## Verification

Fresh verification passed: Python compile, expected-fail current-Go comparison, parser/leak checks, and GitNexus scope review.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/compare_tokenizers.py` | 0 | ✅ pass | 0ms |
| 2 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece python tools/compare_tokenizers.py --mode go-current` | 2 | ✅ expected fail — current Go tokenizer mismatch artifact written | 0ms |
| 3 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece python parser/leak check` | 0 | ✅ pass — s02_artifact_check=pass; raw_probe_text_leaks=0 | 0ms |
| 4 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ⚠️ medium — affected processes are comparator tool flows only, no API runtime flows | 0ms |

## Deviations

GitNexus reported medium risk because `tools/compare_tokenizers.py` changed its own newly indexed script flows; affected processes are internal comparator flows only, not API runtime flows.

## Known Issues

Current Go tokenizer parity still fails, as intended for this slice. GitNexus medium scope is limited to the comparator tool's own flows.

## Files Created/Modified

- `tools/compare_tokenizers.py`
- `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`
- `.gsd/milestones/M012-3edtlz/slices/S02/S02-RESEARCH.md`
