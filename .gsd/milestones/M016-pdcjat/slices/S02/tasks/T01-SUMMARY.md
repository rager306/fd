---
id: T01
parent: S02
milestone: M016-pdcjat
key_files:
  - tools/diagnose_onnx_sequence_length.py
key_decisions:
  - The diagnostic compares TEI API vectors to local Python ONNX outputs at configurable sequence lengths, isolating max-length effects from Go service code.
  - Artifacts include vector hashes and norms, but no raw legal text.
duration: 
verification_result: passed
completed_at: 2026-05-20T05:39:29.333Z
blocker_discovered: false
---

# T01: Implemented the ONNX sequence-length diagnostic script.

**Implemented the ONNX sequence-length diagnostic script.**

## What Happened

Implemented `tools/diagnose_onnx_sequence_length.py`. The script loads M016 S01 worst cases, resolves corpus text by hash, requests TEI vectors from the API, runs the local ONNX artifact with the HF tokenizer at selected sequence lengths, and reports per-length cosine summaries and worst IDs without raw text.

## Verification

`python3 -m py_compile tools/diagnose_onnx_sequence_length.py` passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/diagnose_onnx_sequence_length.py` | 0 | ✅ pass | 0ms |

## Deviations

None.

## Known Issues

This diagnostic depends on TEI API availability and local ONNX Runtime Python packages during T02. It does not test the Go tagged ONNX service directly.

## Files Created/Modified

- `tools/diagnose_onnx_sequence_length.py`
