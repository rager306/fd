---
id: T02
parent: S02
milestone: M012-3edtlz
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizer-go-current-m012-s02.txt
key_decisions:
  - `go-current` mode returns exit code 2 when parity fails, while still writing the mismatch artifact.
  - The current Go tokenizer is confirmed non-equivalent for all five fixed probes.
  - Mismatch windows include nearby token IDs for debugging but omit raw text.
duration: 
verification_result: passed
completed_at: 2026-05-20T02:06:55.747Z
blocker_discovered: false
---

# T02: Compared the current Go tokenizer against the HF baseline and persisted exact mismatch evidence; parity fails for all fixed probes.

**Compared the current Go tokenizer against the HF baseline and persisted exact mismatch evidence; parity fails for all fixed probes.**

## What Happened

Extended `tools/compare_tokenizers.py` with `go-current` mode. It loads the S01 HF baseline, reads the fixed probes via AST, invokes a temporary Go probe under the `api` module using the current `sugarme/tokenizer` path, and writes `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`. The command intentionally exits with code 2 when parity fails. Current Go tokenization failed all five probes: token counts differ and input IDs diverge, confirming the M011 cosine blocker at token level.

## Verification

Python compile passed. `go-current` mode wrote the mismatch artifact and exited with code 2, which is expected because parity failed. The artifact has no raw probe text and records mismatch details by label.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/compare_tokenizers.py && uv run --python 3.13 --with transformers --with torch --with sentencepiece python tools/compare_tokenizers.py --mode go-current` | 2 | ✅ expected fail — artifact written; all 5 probes fail parity | 0ms |
| 2 | `read benchmark-results/fd-tokenizer-go-current-m012-s02.txt` | 0 | ✅ pass — mismatch artifact present with raw_probe_texts_logged=false | 0ms |

## Deviations

The initial implementation had to be corrected so the Go subprocess receives an absolute tokenizer.json path. This avoids cwd-dependent failures when `go run` executes from the `api` module.

## Known Issues

Current `sugarme/tokenizer/pretrained.FromFile + EncodeSingle(true)` fails parity for all fixed probes. Token counts are larger than HF baseline on all probes, and first ID mismatch appears as early as index 2.

## Files Created/Modified

- `tools/compare_tokenizers.py`
- `benchmark-results/fd-tokenizer-go-current-m012-s02.txt`
