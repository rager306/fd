---
id: S01
parent: M005-pyn4x7
milestone: M005-pyn4x7
provides:
  - Accurate benchmark workflow docs for S02/S03 hardening documentation.
requires:
  []
affects:
  - S02
  - S03
key_files:
  - README.md
key_decisions:
  - Document benchmark values as local snapshots rather than universal service guarantees.
  - Make benchmark side effects explicit: Redis FLUSHALL and API restart.
patterns_established:
  - Benchmark docs must state runtime side effects, not only commands.
observability_surfaces:
  - README local benchmark section
drill_down_paths:
  - .gsd/milestones/M005-pyn4x7/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M005-pyn4x7/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M005-pyn4x7/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T10:41:43.319Z
blocker_discovered: false
---

# S01: Benchmark README update

**S01 made README benchmark instructions accurate for uv Python 3.13 and current diagnostic behavior.**

## What Happened

S01 updated README benchmark documentation to match the validated workflow. README now points to benchmark artifacts, uses `uv run --python 3.13 --with requests --with redis`, documents Docker/Redis prerequisites, describes FLUSHALL and API restart side effects, and warns against shared/production benchmark targets. Verification confirmed compose config, benchmark syntax, required snippets, and GitNexus low-risk change detection.

## Verification

All S01 tasks complete and verified.

## Requirements Advanced

- Launchability/operability documentation improved. — 

## Requirements Validated

- README contains uv Python 3.13 benchmark command. — 
- README documents Redis localhost prerequisite and API restart side effect. — 
- docker compose config and benchmark py_compile pass. — 

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

No full benchmark rerun was needed because S01 changed documentation only; compose config and benchmark syntax were verified.

## Known Limitations

Performance numbers remain local-environment snapshots and should be refreshed when hardware/runtime changes.

## Follow-ups

S02 should document Redis overcommit and TEI ONNX/Candle operational notes and record a decision for measured ONNX evaluation.

## Files Created/Modified

- `README.md` — Updated performance and local benchmark instructions for uv Python 3.13, localhost Redis, and benchmark side effects.
