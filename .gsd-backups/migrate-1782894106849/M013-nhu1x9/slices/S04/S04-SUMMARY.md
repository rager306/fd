---
id: S04
parent: M013-nhu1x9
milestone: M013-nhu1x9
provides:
  - Final M013 recommendation.
  - Benchmark-ready tagged ONNX path for future milestone.
requires:
  []
affects:
  []
key_files:
  - .gsd/milestones/M013-nhu1x9/slices/S04/S04-RESEARCH.md
  - benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt
key_decisions:
  - Tagged ONNX path is benchmark-ready on fixed probes.
  - TEI remains production/default.
  - Production switch remains blocked by Docker/CI packaging, native artifact pinning, and larger Russian/legal quality gates.
patterns_established:
  - Correctness gates precede performance gates.
  - Tagged native paths must preserve default build safety.
  - Fixed-probe cosine unlocks benchmarking but not production rollout.
observability_surfaces:
  - S04 research records commands, gates, and next benchmark setup.
  - T02 summary records final verification evidence.
drill_down_paths:
  - .gsd/milestones/M013-nhu1x9/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M013-nhu1x9/slices/S04/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T04:00:51.321Z
blocker_discovered: false
---

# S04: Final benchmark readiness decision

**S04 closed M013 with a benchmark-readiness decision for the tagged ONNX path.**

## What Happened

S04 synthesized the final M013 state and verified it. M013 delivered a validated native artifact contract, safe build-tag boundary, tagged ONNX runtime integration, and fixed-probe cosine equivalence. The final recommendation is to proceed to tagged ONNX performance benchmarking while keeping TEI default and not making production readiness claims.

## Verification

Final verification passed: default tests 78 passed, lint 0 issues, tagged tests 20 passed, health ok, artifact/leak/native checks passed, no background processes, GitNexus low risk.

## Requirements Advanced

- M012-native-packaging-requirement — Validated native artifact/build-tag integration and fixed-probe cosine equivalence for tagged ONNX.

## Requirements Validated

None.

## New Requirements Surfaced

- Need tagged Docker/CI packaging and pinned native artifact release before production use.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None. S04 did not run a performance benchmark; it explicitly recommends a separate benchmark milestone now that tagged cosine passes.

## Known Limitations

No Docker/CI tagged build support yet. No broader corpus quality validation yet. No performance benchmark yet.

## Follow-ups

Plan a tagged ONNX performance benchmarking milestone comparing TEI+Redis vs tagged ONNX+HF tokenizer for cold/warm/batch/cache/startup/memory metrics, with sanitized config snapshots and isolated Redis namespaces.

## Files Created/Modified

- `.gsd/milestones/M013-nhu1x9/slices/S04/S04-RESEARCH.md` — Final M013 benchmark-readiness synthesis.
