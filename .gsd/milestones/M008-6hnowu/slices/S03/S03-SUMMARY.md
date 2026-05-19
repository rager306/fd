---
id: S03
parent: M008-6hnowu
milestone: M008-6hnowu
provides:
  - Final optimization recommendation.
  - Next implementation milestone proposal.
  - Required verification and non-goal list.
requires:
  []
affects:
  []
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S03/S03-RESEARCH.md
key_decisions:
  - Next implementation should be benchmark/config/cache foundation first.
  - Do not adopt ONNX/provider/INT8/language rewrite before measurement and Russian legal quality gates.
  - Reject full C service and model replacement for now.
patterns_established:
  - Boring evidence foundation before risky runtime changes.
  - Explicit non-goals for optimization milestones.
  - Quality gate before model-risking changes.
observability_surfaces:
  - Next milestone must verify sanitized config snapshots, Redis diagnostics, benchmark metrics, and quality-gate readiness.
drill_down_paths:
  - .gsd/milestones/M008-6hnowu/slices/S03/tasks/T01-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T17:14:56.897Z
blocker_discovered: false
---

# S03: Recommend optimization path

**S03 finalized M008: measurement and Redis cache foundation first, ONNX and rewrites later only with evidence.**

## What Happened

S03 synthesized all M008 branches into a final ranked path. It recommends a next milestone centered on benchmark comparability and model-aware long-lived Redis cache, followed by Redis MGET/pipeline A/B only if measured. ONNX FP32 dense-only, provider tuning, INT8, and Rust sidecar are explicitly later gated experiments. This turns the research milestone into an implementation-ready roadmap without speculative runtime migration.

## Verification

S03 research artifact saved and task completed.

## Requirements Advanced

- R001 — Final recommendation keeps Russian legal quality as mandatory gate.
- R002 — Final recommendation prioritizes long-lived Redis cache implementation.
- R003 — Final recommendation includes env-configurable cache/runtime settings.
- R004 — Final recommendation prioritizes sanitized benchmark config snapshots.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S03 originally had no task plan; it was planned and executed after S02 produced the integration/benchmark design.

## Known Limitations

Recommendation does not implement code changes; it intentionally defers runtime changes until evidence infrastructure exists.

## Follow-ups

Create the next implementation milestone for measured cache and benchmark foundation if the user wants execution to continue. Do not push without explicit confirmation.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S03/S03-RESEARCH.md` — Final ranked optimization recommendation.
