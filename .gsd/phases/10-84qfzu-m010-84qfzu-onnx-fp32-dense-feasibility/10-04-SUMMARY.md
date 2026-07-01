---
id: S04
parent: M010-84qfzu
milestone: M010-84qfzu
provides:
  - Final M010 recommendation.
  - Future ONNX adapter gate list.
  - Milestone validation evidence.
requires:
  []
affects:
  []
key_files:
  - .gsd/milestones/M010-84qfzu/slices/S04/S04-RESEARCH.md
  - .gsd/DECISIONS.md
key_decisions:
  - D006: ONNX is locally feasible but not sufficient for production/default runtime switch.
  - Next work should be opt-in adapter/prototype with gates, not a default runtime replacement.
patterns_established:
  - Spike recommendations must separate feasibility from production readiness.
  - Runtime model changes require quality, performance, artifact, and operational gates even after comparator success.
observability_surfaces:
  - S04 research artifact with evidence paths and ONNX hash.
  - D006 decision entry.
  - Final verification task summary with command evidence.
drill_down_paths:
  - .gsd/milestones/M010-84qfzu/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M010-84qfzu/slices/S04/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T18:50:58.645Z
blocker_discovered: false
---

# S04: ONNX spike recommendation

**S04 closed the spike with a bounded recommendation: proceed only to gated opt-in ONNX prototype, not production switch.**

## What Happened

S04 synthesized the completed ONNX spike and ran final verification. It saved a recommendation artifact and decision D006: the exact-model FP32 dense-only ONNX path is locally feasible but should only proceed as a non-default adapter/prototype behind explicit gates. Final verification passed across Go tests, lint, Docker Compose config, Python script compilation/artifact checks, raw-probe leakage checks, and GitNexus change detection. Production runtime remains unchanged.

## Verification

S04 verification passed with fresh evidence in T02: Go tests 60 passed, GolangCI-Lint 0 issues, Docker Compose config OK, Python py_compile/artifact checks OK, no raw probe leakage, GitNexus low risk/no affected processes.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

- Future ONNX adapter milestone should define artifact distribution/checksum handling for large model files.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None. S04 remained a recommendation/verification slice and did not alter production runtime defaults.

## Known Limitations

No throughput benchmark, larger retrieval-quality benchmark, artifact distribution plan, or Go API ONNX adapter was built in M010. Those are intentionally future gates.

## Follow-ups

Plan a future non-default ONNX adapter/prototype milestone only if desired. It should start with artifact distribution/checksum design, then opt-in backend config, performance benchmarking, and larger Russian/legal quality gate.

## Files Created/Modified

- `.gsd/milestones/M010-84qfzu/slices/S04/S04-RESEARCH.md` — Final ONNX spike recommendation and future gates.
- `.gsd/DECISIONS.md` — Durable runtime decision D006.
- `.gsd/milestones/M010-84qfzu/slices/S04/tasks/T02-SUMMARY.md` — Final verification task summary.
