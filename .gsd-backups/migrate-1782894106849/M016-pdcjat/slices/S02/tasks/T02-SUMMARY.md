---
id: T02
parent: S02
milestone: M016-pdcjat
key_files:
  - tools/diagnose_onnx_sequence_length.py
  - benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt
key_decisions:
  - ONNX sequence length 128 is confirmed as the primary root cause for the severe M015 vector divergence on the profiled worst cases.
  - A 512-token ONNX run nearly restores TEI equivalence for all worst cases that fit within 512; two cases still truncate at 512 and remain around 0.990 cosine.
duration: 
verification_result: passed
completed_at: 2026-05-20T05:41:08.941Z
blocker_discovered: false
---

# T02: Ran sequence-length diagnostics and confirmed 128-token truncation as the main M015 divergence cause.

**Ran sequence-length diagnostics and confirmed 128-token truncation as the main M015 divergence cause.**

## What Happened

Ran the sequence-length diagnostic against TEI API and local ONNX outputs for the 17 worst M015 cases. At max length 128, all 17 cases are truncated and mean cosine is 0.9205, with minimum 0.7439. At max length 512, only 2 cases remain truncated and mean cosine rises to 0.998856, with non-truncated worst cases near 0.99999. This confirms truncation at 128 as the primary cause for severe divergence, while very long cases may still need chunking or longer sequence support.

## Verification

Diagnostic run passed, artifact includes 128 vs 512 summaries, and raw legal text leak check passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests --with transformers --with torch --with sentencepiece --with numpy --with onnxruntime python tools/diagnose_onnx_sequence_length.py ... --sequence-lengths 128,512` | 0 | ✅ pass — diagnostic artifact written | 34100ms |
| 2 | `python sequence diagnostic artifact hygiene check` | 0 | ✅ pass — sequence_diagnostics_hygiene=pass; raw_legal_text_leaks=0 | 0ms |

## Deviations

None.

## Known Issues

The diagnostic uses Python ONNX Runtime directly, not the Go tagged ONNX service. S03 should decide whether to re-export and wire a 512 max sequence artifact or implement chunking; performance impact remains unknown.

## Files Created/Modified

- `tools/diagnose_onnx_sequence_length.py`
- `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt`
