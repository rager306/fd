---
id: S02
parent: M016-pdcjat
milestone: M016-pdcjat
provides:
  - Root-cause evidence for S03 remediation planning.
  - Sanitized TEI-vs-ONNX 128/512 cosine artifact.
requires:
  []
affects:
  - S03
key_files:
  - tools/diagnose_onnx_sequence_length.py
  - benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt
key_decisions:
  - Primary cause confirmed: ONNX max sequence length 128 truncation.
  - 512-token ONNX fixes most worst cases, but the two cases still truncated at 512 remain below strict 0.999 cosine.
  - Do not prioritize model alternatives or packaging until longer-sequence/chunking remediation is tested.
patterns_established:
  - Use local ONNX sequence-length diagnostics before re-exporting or changing the service.
  - Quality remediation must account for both model max length and legal chunking policy.
observability_surfaces:
  - Sequence diagnostic artifact with per-length cosine summaries, worst IDs, hashes, norms, truncation flags, and vector hashes.
drill_down_paths:
  - .gsd/milestones/M016-pdcjat/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M016-pdcjat/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M016-pdcjat/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T07:03:46.115Z
blocker_discovered: false
---

# S02: Sequence length root cause diagnostics

**S02 confirmed ONNX 128-token truncation as the root cause of the M015 severe legal divergence.**

## What Happened

S02 implemented and ran focused sequence-length diagnostics. The results are decisive: current ONNX 128-token truncation explains the severe legal vector divergence. Increasing local ONNX max length to 512 raises mean cosine on the 17 worst cases from 0.9204953 to 0.99885631 and restores near-0.99999 cosine for cases that fit within 512. Two longer cases remain truncated at 512 and still fall to about 0.990, so chunking or longer max length remains necessary for complete legal quality.

## Verification

S02 verification passed: `python3 -m py_compile tools/diagnose_onnx_sequence_length.py`; TEI health ok; diagnostic run completed; artifact hygiene check passed with `raw_legal_text_leaks=0`.

## Requirements Advanced

- onnx-long-text-quality — Confirmed the severe legal divergence failure mode and narrowed remediation to sequence length/chunking.

## Requirements Validated

- m016-s02-root-cause — `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt` shows 128 mean cosine 0.9204953 vs 512 mean cosine 0.99885631 on the same worst IDs.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Первый вызов закрытия S02 был отклонён валидатором из-за отсутствующих обязательных полей; данные задач не потеряны, повторяю корректный вызов.

## Known Limitations

The diagnostic used local Python ONNX Runtime, not the Go tagged ONNX service. A 512-token Go/runtime artifact still needs implementation and a full corpus gate rerun.

## Follow-ups

S03 should choose remediation: re-export/wire ONNX max sequence length 512 and rerun the legal gate, plus define chunking or longer sequence strategy for texts over 512 tokens. Performance impact must be measured only after quality passes.

## Files Created/Modified

- `tools/diagnose_onnx_sequence_length.py` — Focused TEI-vs-local-ONNX sequence length diagnostic tool.
- `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt` — Diagnostic artifact comparing ONNX max length 128 vs 512 on M015 worst cases.
