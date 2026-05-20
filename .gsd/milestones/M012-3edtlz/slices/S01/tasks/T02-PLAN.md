---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Generate Hugging Face tokenizer baseline

Add a Python tool that loads the local Hugging Face tokenizer for `deepvk/USER-bge-m3`, reuses the existing fixed probes, and writes a sanitized baseline artifact under benchmark-results.

## Inputs

- `tools/compare_dense_embeddings.py`
- `tei-models/deepvk--USER-bge-m3/tokenizer.json`

## Expected Output

- `tools/compare_tokenizers.py`
- `benchmark-results/fd-tokenizer-baseline-m012-s01.txt`

## Verification

`uv run --python 3.13 --with transformers --with torch --with sentencepiece python tools/compare_tokenizers.py --mode baseline` writes the artifact and exits 0.

## Observability Impact

Creates deterministic expected tokenizer evidence for future Go comparisons.
