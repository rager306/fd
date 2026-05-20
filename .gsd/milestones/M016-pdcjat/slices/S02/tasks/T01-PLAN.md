---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Implement ONNX sequence length diagnostic

Create a focused diagnostic script that loads worst M015 IDs, requests TEI vectors, runs the local ONNX artifact through HF tokenizer at configurable max sequence lengths, and writes sanitized cosine results without raw legal text.

## Inputs

- `tools/profile_legal_divergence.py`
- `tools/compare_dense_embeddings.py`
- `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`

## Expected Output

- `tools/diagnose_onnx_sequence_length.py`

## Verification

`python3 -m py_compile tools/diagnose_onnx_sequence_length.py` passes.

## Observability Impact

Creates repeatable sequence-length diagnostics for the suspected root cause.
