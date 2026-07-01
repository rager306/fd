---
id: S05
parent: M008-6hnowu
milestone: M008-6hnowu
provides:
  - Language/runtime rewrite strategy for final M008 recommendation.
  - Stop criteria for rewrite experiments.
requires:
  []
affects:
  []
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S05/S05-RESEARCH.md
key_decisions:
  - Keep Go service for now.
  - Do not full-rewrite in Rust or C without pprof/per-layer evidence.
  - Rust sidecar is the only plausible rewrite experiment; C is limited to narrow FFI if needed.
patterns_established:
  - Measured bottleneck before language rewrite.
  - Sidecar A/B beats full rewrite for risky runtime experiments.
  - C only for narrow FFI when wrappers are insufficient.
observability_surfaces:
  - Future rewrite gate requires pprof/per-layer timing, cache/Redis/model latency decomposition, and sanitized config snapshots.
drill_down_paths:
  - .gsd/milestones/M008-6hnowu/slices/S05/tasks/T01-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S05/tasks/T02-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S05/tasks/T03-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S05/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T17:11:43.348Z
blocker_discovered: false
---

# S05: Research Go vs C vs Rust performance tradeoffs

**S05 closed the rewrite question: no full rewrite now; measure first, keep Go, maybe Rust sidecar later.**

## What Happened

S05 mapped fd's language-sensitive bottlenecks, evaluated Rust and C ecosystem maturity, and produced a rewrite strategy. The conclusion is conservative: optimize and measure within Go first, prioritize Redis and ONNX/provider/threading evidence, consider Rust only as an A/B sidecar if native inference makes it worthwhile, and avoid C except as a small audited boundary.

## Verification

All S05 tasks complete and S05 research artifact saved.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

No pprof profile was collected in this research slice; conclusion is based on architecture mapping and ecosystem/source evidence.

## Follow-ups

Add pprof/per-layer timing in a future implementation milestone before revisiting rewrite decisions. If ONNX native inference becomes a real path, compare Go adapter and Rust sidecar under the same Russian/legal quality gate.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S05/S05-RESEARCH.md` — Language/runtime tradeoff research and recommendation artifact.
