---
id: S01
parent: M014-vjfs9f
milestone: M014-vjfs9f
provides:
  - Benchmark metadata harness for S02/S03.
requires:
  []
affects:
  - S02
  - S03
key_files:
  - benchmark.py
key_decisions:
  - Use snapshot_version 2 for M014 benchmark artifacts.
  - Record runtime label, build tags, ONNX manifest, native tokenizer manifest, ONNX Runtime library path/hash, and Redis namespace.
  - Preserve existing TEI/default benchmark behavior when env vars are unset.
patterns_established:
  - Benchmark metadata changes precede expensive runtime measurements.
  - Optional metadata preserves default benchmark compatibility.
  - Runtime artifacts must record build mode and native/ONNX/ORT checksums.
observability_surfaces:
  - Benchmark snapshot v2 runtime metadata.
  - Task summaries documenting matrix and required env vars.
drill_down_paths:
  - .gsd/milestones/M014-vjfs9f/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M014-vjfs9f/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M014-vjfs9f/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T04:13:56.252Z
blocker_discovered: false
---

# S01: Benchmark matrix and metadata harness

**S01 prepared benchmark snapshot v2 so TEI and tagged ONNX runs can be compared reproducibly.**

## What Happened

S01 defined the benchmark matrix and extended the harness metadata. `benchmark.py` now records optional tagged ONNX runtime metadata in the sanitized effective configuration snapshot while preserving default behavior. Dry-run checks confirmed snapshot v2 includes tagged build/native/ONNX/ORT fields and excludes raw probe text. This makes subsequent TEI and tagged ONNX artifacts comparable.

## Verification

py_compile passed; snapshot dry-run passed with tagged ONNX metadata; raw probe leaks 0; GitNexus medium scope reviewed as benchmark-only.

## Requirements Advanced

- benchmark-comparability — Added runtime/native/ONNX/ORT metadata fields required for benchmark comparability.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S01 did not run full runtime benchmarks; it intentionally stopped at metadata harness changes and dry-run checks so subsequent TEI/ONNX runs are comparable.

## Known Limitations

Snapshot records artifact checksums but does not enforce them; artifact validators remain responsible for enforcement.

## Follow-ups

S02 should run a fresh TEI baseline using snapshot v2. S3 should run tagged ONNX with env metadata fields populated and isolated Redis namespace.

## Files Created/Modified

- `benchmark.py` — Benchmark config snapshot v2 with optional tagged ONNX/native/ORT metadata.
