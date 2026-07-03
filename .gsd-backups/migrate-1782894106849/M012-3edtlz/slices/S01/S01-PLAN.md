# S01: Hugging Face tokenizer baseline

**Goal:** Create a deterministic, sanitized Hugging Face tokenizer baseline for the fixed Russian/legal probe set used by the embedding comparators.
**Demo:** After this, the project has a sanitized Hugging Face tokenizer baseline for fixed Russian/legal probes.

## Must-Haves

- Uses local `tei-models/deepvk--USER-bge-m3` tokenizer.
- Reuses fixed probe labels without printing raw probe texts.
- Captures `input_ids` and `attention_mask` evidence with hashes and lengths.
- Records tokenizer/model revision and relevant package versions.
- Artifact parser and raw-text leakage checks pass.

## Proof Level

- This slice proves: Executable Python script plus generated artifact and leak/parser checks.

## Integration Closure

Downstream S02 will compare Go tokenizer candidates against this baseline instead of guessing from embedding cosine differences.

## Verification

- Adds safe tokenization evidence: labels, char counts, token lengths, masks, hashes, tokenizer/model metadata, and raw-text exclusion.

## Tasks

- [x] **T01: Design tokenizer baseline artifact format** `est:small`
  Inspect the existing dense comparator probe structure and design a tokenizer-baseline artifact format that records safe metadata, token IDs or token hashes, and raw-text exclusion policy.
  - Files: `tools/compare_dense_embeddings.py`
  - Verify: Artifact schema is documented in task summary and excludes raw probe text.

- [x] **T02: Generate Hugging Face tokenizer baseline** `est:medium`
  Add a Python tool that loads the local Hugging Face tokenizer for `deepvk/USER-bge-m3`, reuses the existing fixed probes, and writes a sanitized baseline artifact under benchmark-results.
  - Files: `tools/compare_tokenizers.py`, `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`
  - Verify: `uv run --python 3.13 --with transformers --with torch --with sentencepiece python tools/compare_tokenizers.py --mode baseline` writes the artifact and exits 0.

- [x] **T03: Verify baseline artifact hygiene** `est:small`
  Verify the tokenizer baseline artifact: parse it, check expected sections, ensure raw probe texts are absent, run Python compile, and record evidence.
  - Files: `tools/compare_tokenizers.py`, `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`
  - Verify: Python compile, artifact parser, raw-text leakage check, and GitNexus detect_changes pass.

## Files Likely Touched

- tools/compare_dense_embeddings.py
- tools/compare_tokenizers.py
- benchmark-results/fd-tokenizer-baseline-m012-s01.txt
