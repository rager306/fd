# S01: Hugging Face tokenizer baseline — UAT

**Milestone:** M012-3edtlz
**Written:** 2026-05-20T02:00:43.637Z

# S01 UAT — Hugging Face tokenizer baseline

## Checks

- [x] Baseline tool exists at `tools/compare_tokenizers.py`.
- [x] Tool loads local `tei-models/deepvk--USER-bge-m3` tokenizer.
- [x] Tool reuses existing fixed probes without rendering raw text.
- [x] Baseline artifact exists at `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`.
- [x] Artifact records tokenizer metadata, token IDs/masks, hashes, and `raw_probe_texts_logged=false`.
- [x] Parser/leak check passed with `raw_probe_text_leaks=0`.
- [x] GitNexus reported low risk and no affected processes.

## UAT Result

Pass. S02 can now compare Go tokenizer outputs against the HF baseline.

