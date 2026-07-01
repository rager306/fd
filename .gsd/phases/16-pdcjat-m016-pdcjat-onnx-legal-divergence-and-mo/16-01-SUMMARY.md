---
id: S01
parent: M016-pdcjat
milestone: M016-pdcjat
provides:
  - Worst-case token/truncation target set for S02 runtime diagnostics.
requires:
  []
affects:
  - S02
  - S03
key_files:
  - tools/profile_legal_divergence.py
  - benchmark-results/fd-legal-divergence-profile-m016-s01.txt
key_decisions:
  - Max sequence length 128 is now the primary suspect because all 17 worst cases truncate at 128, while only 2 truncate at 512.
  - S02 should test 512-token behavior before larger remediation work.
patterns_established:
  - Long-text quality failures should be profiled by stable IDs, hashes, and token counts before runtime changes.
  - No raw legal corpus text in diagnostic artifacts.
observability_surfaces:
  - Divergence profile artifact with IDs, hashes, token counts, truncation flags, and tokenizer hash.
drill_down_paths:
  - .gsd/milestones/M016-pdcjat/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M016-pdcjat/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M016-pdcjat/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T05:35:58.175Z
blocker_discovered: false
---

# S01: Worst-case divergence profile

**S01 proved M015 worst legal divergence cases are all truncated by the current 128-token ONNX path.**

## What Happened

S01 converted the M015 legal quality failure into a concrete diagnostic set. It extracted worst IDs, resolved them against the 44-ФЗ corpus, verified hashes, and profiled token lengths using the local Hugging Face tokenizer. All 17 worst cases exceed 128 tokens; only 2 exceed 512. This strongly points to the current ONNX max sequence length 128 path as the likely cause of severe vector divergence.

## Verification

All S01 verification checks passed.

## Requirements Advanced

- onnx-long-text-quality — Narrowed legal quality divergence to long-text truncation risk.

## Requirements Validated

- m015-worst-case-profile — `benchmark-results/fd-legal-divergence-profile-m016-s01.txt` shows 17/17 worst cases truncated at 128 and all hashes match.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Tokenizer diagnostics identify truncation risk but do not directly expose TEI internal behavior. Runtime confirmation is required in S02.

## Follow-ups

S02 should run focused runtime diagnostics on these IDs: compare TEI output, current tagged ONNX 128 output, and a local Python/ONNX or export variant at 512 where feasible. The main hypothesis is max sequence length 128 truncation.

## Files Created/Modified

- `tools/profile_legal_divergence.py` — Sanitized profiler for M015 worst divergence IDs.
- `benchmark-results/fd-legal-divergence-profile-m016-s01.txt` — Worst-case token/truncation profile artifact.
