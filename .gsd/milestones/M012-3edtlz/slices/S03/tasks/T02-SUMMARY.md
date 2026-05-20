---
id: T02
parent: S03
milestone: M012-3edtlz
key_files:
  - tools/compare_tokenizers.py
  - benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt
key_decisions:
  - `github.com/daulet/tokenizers` with HF Rust `libtokenizers.a` passes token-level parity for all five fixed probes.
  - Candidate comparison remains isolated from project runtime dependencies for now.
  - S03 T03 can proceed to integration if packaging constraints are accepted as prototype-only.
duration: 
verification_result: passed
completed_at: 2026-05-20T02:19:20.287Z
blocker_discovered: false
---

# T02: Compared the HF Rust tokenizers Go binding against the baseline and proved token-level parity for all fixed probes.

**Compared the HF Rust tokenizers Go binding against the baseline and proved token-level parity for all fixed probes.**

## What Happened

Added `go-hf-binding` mode to the tokenizer comparator. It loads the S01 baseline, creates a temporary Go module, imports `github.com/daulet/tokenizers`, links against a provided `libtokenizers.a`, loads the local `tokenizer.json`, and compares all fixed probes against the HF baseline. The comparison passed for all five probes: token counts, input ID hashes, and attention mask hashes match exactly. The artifact is sanitized and records `raw_probe_texts_logged=false`.

## Verification

Python compile passed. `go-hf-binding` mode exited 0 and wrote `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt` with PASS and no raw probe text.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/compare_tokenizers.py && uv run --python 3.13 --with transformers --with torch --with sentencepiece python tools/compare_tokenizers.py --mode go-hf-binding --hf-tokenizers-lib-dir /tmp/fd-daulet-tokenizers-probe` | 0 | ✅ pass — all 5 fixed probes match HF baseline | 0ms |
| 2 | `grep title/raw_probe/PASS benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt` | 0 | ✅ pass — S03 artifact title, raw_probe_texts_logged=false, PASS | 0ms |

## Deviations

Extended `tools/compare_tokenizers.py` with `go-hf-binding` mode to persist the candidate comparison artifact. The mode uses a temp module so the project `api/go.mod` is not modified until integration is explicitly chosen.

## Known Issues

The passing comparison depends on a local temp prebuilt `libtokenizers.a` under `/tmp/fd-daulet-tokenizers-probe`. Production/CI packaging remains unresolved and must be documented if integrated.

## Files Created/Modified

- `tools/compare_tokenizers.py`
- `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt`
