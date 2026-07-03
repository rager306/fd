---
id: S03
parent: M016-pdcjat
milestone: M016-pdcjat
provides:
  - Remediation plan for future implementation milestone.
  - D014 project decision.
requires:
  []
affects:
  - Future ONNX remediation milestone
key_files:
  - benchmark-results/fd-onnx-remediation-plan-m016-s03.txt
  - .gsd/DECISIONS.md
key_decisions:
  - D014: keep TEI production/default; keep ONNX opt-in; next ONNX gate is 512-token quality remediation plus long-text policy before packaging/tuning/promotion.
  - Reject model switch or production ONNX promotion based on M016 evidence.
patterns_established:
  - Do not promote runtime optimizations until legal-quality gates pass.
  - Long legal text requires explicit sequence-length/chunking policy; tokenizer parity alone is insufficient.
observability_surfaces:
  - Remediation plan artifact with metrics, option tradeoffs, acceptance gates, and no raw legal text.
drill_down_paths:
  - .gsd/milestones/M016-pdcjat/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M016-pdcjat/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M016-pdcjat/slices/S03/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T07:20:06.153Z
blocker_discovered: false
---

# S03: Remediation path decision

**S03 chose the ONNX legal-divergence remediation path and recorded the decision.**

## What Happened

S03 completed the remediation decision work. It used M015/S01/S02 measurements to compare remediation options and selected a quality-first path: keep TEI as the default, keep ONNX experimental, validate a 512-token ONNX path first, and add chunking or longer-sequence handling for legal texts beyond 512 tokens. The slice also recorded D014 so future work does not package or promote the known-bad 128-token path.

## Verification

Fresh verification passed: Python compile and artifact hygiene passed; Go short tests passed with 78 tests in 4 packages; pinned GolangCI-Lint reported 0 issues; tagged HF tokenizer tests passed with 20 tests in 1 package; GitNexus scope check was low.

## Requirements Advanced

- onnx-long-text-quality — Converted root-cause evidence into a concrete quality-first remediation path.

## Requirements Validated

- m016-remediation-decision — `benchmark-results/fd-onnx-remediation-plan-m016-s03.txt` and D014 record the remediation path and acceptance gates.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

No runtime remediation was implemented in S03. Full legal retrieval gate rerun remains future work after a 512-token/chunking implementation exists.

## Follow-ups

Create a future milestone to implement and validate the 512-token ONNX quality gate, then define chunking/longer-sequence handling for >512-token legal fragments and rerun full legal retrieval quality before performance or packaging work.

## Files Created/Modified

- `benchmark-results/fd-onnx-remediation-plan-m016-s03.txt` — Remediation plan comparing 512-token ONNX, chunking, longer sequence, and rejected model/default switch options.
- `.gsd/DECISIONS.md` — Decision register updated with D014 remediation path.
