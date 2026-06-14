---
id: S01
parent: M042-fjf2en
milestone: M042-fjf2en
provides:
  - RCA artifact for S02 implementation choices.
  - Evidence-backed TEI-first direction and ONNX deferral.
requires:
  []
affects:
  []
key_files:
  - documents/te-perf-snapshot-m042-s01.md
  - benchmark-results/te-concurrency-profile-m042-s01.md
  - documents/te-perf-root-cause-m042.md
key_decisions:
  - D047: ONNX is removed from current runtime strategy and kept only as future research.
  - R022 deferred; R020 validated.
patterns_established:
  - Do not use destructive restart profiling as a casual perf measurement for TEI.
  - Separate fd cache-hot acceptance from direct TEI miss/runtime diagnostics.
observability_surfaces:
  - TEI `/info`, docker health logs, and TEI server-side timing logs were used as evidence.
drill_down_paths:
  - .gsd/milestones/M042-fjf2en/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M042-fjf2en/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M042-fjf2en/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-06-14T10:20:03.157Z
blocker_discovered: false
---

# S01: TEI queue_time root cause analysis

**Completed TEI RCA showing batch-size-sensitive queue_time and severe backend startup delay, with ONNX deferred from M042.**

## What Happened

S01 captured direct TEI telemetry, documented a failed/destructive concurrency profile attempt, and wrote the final TEI RCA. T01 proved direct TEI batch-size-sensitive queue_time: batch=1 queue p50 ~0.253ms, batch=8 ~427ms, batch=32 ~2435ms, bypassing fd and Redis. T02 documented that restart/recreate profiling is operationally risky: TEI spent roughly 48 minutes from `Starting model backend` to `Ready`, including delayed missing-ONNX ORT fallback before Candle/safetensors warmup. T03 synthesized the RCA and rejected ONNX as a current M042 mitigation; ONNX is now future research only, while M042 proceeds TEI-first.

## Verification

T01/T02/T03 verification passed, and structured runtime UAT was saved with four gsd_uat_exec evidence checks. R020 was validated.

## Requirements Advanced

- R020 — RCA evidence and verdict completed.
- R027 — Created TEI-first/ONNX-disabled project constraint for follow-up implementation.

## Requirements Validated

- R020 — `documents/te-perf-root-cause-m042.md` validates TEI queue/startup RCA.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Original S01 T02 intended complete parallel concurrency metrics and a restart profile. The restart scenario became destructive, logs were lost after recreate, and the slice was replanned to document this as evidence rather than repeat unsafe restarts. The original M042 ONNX direction was rescoped by user decision and RCA evidence.

## Known Limitations

Parallel TEI concurrency metrics remain incomplete. TEI startup/recreate remains slow; avoid unnecessary restarts until S02 clarifies mitigation. ONNX branch deactivation is not yet implemented; it is now captured as R027 and M042 follow-up scope.

## Follow-ups

Proceed to TEI-first S02: remove/disable active ONNX surfaces and then evaluate any TEI request-shaping/async mitigation with safe evidence.

## Files Created/Modified

None.
