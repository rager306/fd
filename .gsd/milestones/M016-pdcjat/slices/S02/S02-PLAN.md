# S02: Truncation root-cause diagnostics

**Goal:** Run focused TEI/HF/ONNX diagnostics on worst cases to isolate truncation/tokenization/export behavior.
**Demo:** After this, we know whether max sequence length 128 is the likely cause or whether another ONNX path issue remains.

## Must-Haves

- TEI vs current ONNX 128 worst-case reproduction is rerun or referenced.
- Python ONNX sequence-length diagnostics compare 128 and 512 where feasible.
- HF tokenizer token counts and truncation are connected to cosine outcomes.
- Evidence supports or rejects max_sequence_length=128 as primary cause.
- No production switch occurs.

## Proof Level

- This slice proves: Runtime diagnostics and sanitized artifacts.

## Integration Closure

Determines whether to pursue longer ONNX export, chunking, or tokenizer/export fixes.

## Verification

- Produces focused cosine/tokenization artifacts for worst IDs and sequence-length variants.

## Tasks

- [x] **T01: Implement ONNX sequence length diagnostic** `est:medium`
  Create a focused diagnostic script that loads worst M015 IDs, requests TEI vectors, runs the local ONNX artifact through HF tokenizer at configurable max sequence lengths, and writes sanitized cosine results without raw legal text.
  - Files: `tools/diagnose_onnx_sequence_length.py`
  - Verify: `python3 -m py_compile tools/diagnose_onnx_sequence_length.py` passes.

- [x] **T02: Run sequence length diagnostics** `est:medium`
  Run the diagnostic against TEI and local ONNX artifact for sequence lengths 128 and 512, writing `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt`.
  - Files: `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt`
  - Verify: Diagnostic artifact exists and includes 128 vs 512 cosine summaries and no raw legal text.

- [x] **T03: Record root-cause verdict** `est:small`
  Interpret the diagnostic result and decide whether max sequence length 128 is confirmed as root cause or if pooling/export/TEI behavior remains suspect.
  - Verify: Task summary states confirmed/rejected/narrowed cause and next remediation path.

## Files Likely Touched

- tools/diagnose_onnx_sequence_length.py
- benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt
