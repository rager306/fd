---
id: S04
parent: M015-22msl0
milestone: M015-22msl0
provides:
  - Blocking quality decision and next-step recommendation.
requires:
  []
affects:
  []
key_files:
  - benchmark-results/fd-legal-retrieval-m015-summary.txt
  - .gsd/DECISIONS.md
key_decisions:
  - D012: block ONNX packaging/tuning as next priority; investigate long-text legal divergence first.
patterns_established:
  - A failed quality gate can complete a milestone successfully when the milestone objective is to produce the verdict.
  - Packaging/tuning must wait when quality gates expose severe vector equivalence outliers.
observability_surfaces:
  - Summary artifact with key metrics and decision impact.
  - D012 decision.
  - Final verification task summary.
drill_down_paths:
  - .gsd/milestones/M015-22msl0/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M015-22msl0/slices/S04/tasks/T02-SUMMARY.md
  - .gsd/milestones/M015-22msl0/slices/S04/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T05:07:48.570Z
blocker_discovered: false
---

# S04: Quality verdict and closure

**S04 recorded the failed legal quality gate and closed the slice with verification.**

## What Happened

S04 summarized the quality gate, recorded the blocking decision, and ran final verification. The result is clear: tagged ONNX does not pass the Russian/legal quality gate on the provided 44-ФЗ corpus. TEI remains production/default. The next priority is long-text/truncation divergence investigation, not packaging or tuning.

## Verification

All final verification gates passed with fresh output after the last changes.

## Requirements Advanced

None.

## Requirements Validated

- quality-gate-executed — M015 produced corpus profile, evaluator, live FAIL artifact, summary, and final verification.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

- onnx-production-quality — Tagged ONNX failed strict legal-corpus vector equivalence on the 44-ФЗ sample.

## Operational Readiness

None.

## Deviations

The live quality artifact is FAIL by design; S04 completed because the milestone goal is to produce a gate verdict, not force ONNX acceptance.

## Known Limitations

The quality corpus lacks explicit qrels. The gate is parity/known-item evidence, but the vector equivalence outliers are severe enough to block ONNX readiness.

## Follow-ups

Recommended next milestone: investigate ONNX long-text legal divergence. Scope: compare TEI tokenization/truncation behavior against ONNX max_sequence_length=128, test longer sequence ONNX export feasibility, evaluate chunking policy, and rerun the legal gate.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m015-summary.txt` — Quality gate summary with FAIL verdict and next-step recommendation.
- `.gsd/DECISIONS.md` — Decision D012 blocks ONNX packaging/tuning until long-text legal quality divergence is investigated.
