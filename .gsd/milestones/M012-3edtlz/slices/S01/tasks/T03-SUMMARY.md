---
id: T03
parent: S01
milestone: M012-3edtlz
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizer-baseline-m012-s01.txt
key_decisions:
  - Parser/leak checks for tokenizer tooling should run under `uv run --python 3.13 --with transformers --with torch --with sentencepiece` unless the script is refactored to separate pure parsing helpers from tokenizer imports.
duration: 
verification_result: passed
completed_at: 2026-05-20T02:00:09.561Z
blocker_discovered: false
---

# T03: Verified the tokenizer baseline artifact parses and does not leak raw probe text.

**Verified the tokenizer baseline artifact parses and does not leak raw probe text.**

## What Happened

Verified the tokenizer baseline artifact. Python bytecode compilation passed. The artifact contains the expected sections and configuration block, records `raw_probe_texts_logged=false`, has probe_count 5, and contains none of the raw probe strings from the fixed probe source. GitNexus change detection reported low risk with no affected processes; the only code change is a new Python tool, so no existing Go execution flow is affected.

## Verification

Fresh verification passed after generation: Python compile, uv-context parser/leak check, GitNexus detect_changes, and cleanup of accidental `__pycache__`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/compare_tokenizers.py` | 0 | ✅ pass | 0ms |
| 2 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece python artifact parser/leak check` | 0 | ✅ pass — tokenizer_baseline_artifact_check=pass; raw_probe_text_leaks=0 | 0ms |
| 3 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ✅ pass — low risk; affected_processes=[] | 0ms |
| 4 | `rm -rf tools/__pycache__ && git status --short` | 0 | ✅ pass — no pycache remains | 0ms |

## Deviations

The first parser/leak check was run outside the `uv` dependency context and failed to import `transformers`; the check was rerun under the same `uv` dependency context used by the tool and passed.

## Known Issues

None.

## Files Created/Modified

- `tools/compare_tokenizers.py`
- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`
