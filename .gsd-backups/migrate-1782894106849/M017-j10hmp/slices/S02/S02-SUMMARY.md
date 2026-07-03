---
id: S02
parent: M017-j10hmp
milestone: M017-j10hmp
provides:
  - Next remediation recommendation for future GSD milestone.
requires:
  []
affects:
  - Future chunking or longer sequence milestone
key_files:
  - benchmark-results/fd-onnx-512-outcome-m017-s02.txt
  - .gsd/DECISIONS.md
key_decisions:
  - D015: 512-token ONNX is necessary but insufficient; next remediation must handle >512-token legal fragments.
  - Do not promote ONNX based on excellent ranking parity alone because strict vector equivalence still fails.
patterns_established:
  - Ranking parity and vector equivalence are separate gates.
  - Measured FAIL outcomes are valid milestone outputs when they answer the risk question.
observability_surfaces:
  - Outcome artifact comparing M015 128, M016 Python 512, and M017 Go 512 evidence.
drill_down_paths:
  - .gsd/milestones/M017-j10hmp/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M017-j10hmp/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M017-j10hmp/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T07:31:34.958Z
blocker_discovered: false
---

# S02: Quality outcome decision

**S02 recorded that ONNX 512 is necessary but insufficient and chose chunking or longer-sequence handling as the next remediation.**

## What Happened

S02 assessed the measured 512-token gate. It concluded that the tagged Go ONNX 512 path is a major improvement over the 128-token path but does not pass strict legal equivalence. The slice recorded D015 to prevent treating ranking parity as production readiness. The recommended next remediation is hybrid: keep 512 as the experimental baseline and add deterministic chunking or longer-sequence handling for legal fragments above 512 tokens.

## Verification

Fresh verification passed: Python compile and artifact hygiene, Go short tests, pinned lint, tagged HF tokenizer tests, no background processes, and GitNexus scope check.

## Requirements Advanced

- onnx-long-text-quality — Converted ONNX 512 strict FAIL into a concrete next remediation requirement for >512-token legal fragments.

## Requirements Validated

- m017-512-outcome-decision — `benchmark-results/fd-onnx-512-outcome-m017-s02.txt` and D015 state that 512 is necessary but insufficient.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

No chunking or longer-sequence implementation was added in M017. The strict cosine gate remains failed for ONNX 512.

## Follow-ups

Plan a future milestone for deterministic legal chunking or a longer-sequence ONNX variant. The recommended path is hybrid: 512-token ONNX baseline plus chunking for fragments above 512 tokens, with cache namespace versioning and full legal retrieval rerun.

## Files Created/Modified

- `benchmark-results/fd-onnx-512-outcome-m017-s02.txt` — Outcome assessment for ONNX 512 legal quality gate.
- `.gsd/DECISIONS.md` — Decision register updated with D015.
