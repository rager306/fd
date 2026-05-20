---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Run sequence length diagnostics

Run the diagnostic against TEI and local ONNX artifact for sequence lengths 128 and 512, writing `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt`.

## Inputs

- `tools/diagnose_onnx_sequence_length.py`
- `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`
- `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`

## Expected Output

- `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt`

## Verification

Diagnostic artifact exists and includes 128 vs 512 cosine summaries and no raw legal text.

## Observability Impact

Shows whether 512-token ONNX improves worst-case cosine against TEI.
