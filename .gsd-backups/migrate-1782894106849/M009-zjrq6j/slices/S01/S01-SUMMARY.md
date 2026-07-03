---
id: S01
parent: M009-zjrq6j
milestone: M009-zjrq6j
provides:
  - Sanitized benchmark config snapshot foundation for S02/S03/S04.
  - M009 S01 benchmark artifact.
requires:
  []
affects:
  []
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m009-s01.txt
key_decisions:
  - Use an allowlisted environment snapshot rather than dumping process environment.
  - Record Docker compose config and images with hashes/identifiers but avoid printing full compose config.
  - Keep raw benchmark input texts out of artifacts; use labels and lengths instead.
patterns_established:
  - Allowlist env snapshots for committed benchmark artifacts.
  - Print labels/lengths instead of raw benchmark texts.
  - Gracefully degrade metadata collection if Docker/Git/Redis are unavailable.
observability_surfaces:
  - Benchmark snapshot section, Redis INFO before-run summary, Docker compose hash/images, environment baseline hash, git dirty state, redaction policy.
drill_down_paths:
  - .gsd/milestones/M009-zjrq6j/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M009-zjrq6j/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M009-zjrq6j/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T17:30:39.681Z
blocker_discovered: false
---

# S01: Benchmark config snapshot

**S01 made benchmark results comparable by adding a sanitized effective config snapshot to every run.**

## What Happened

S01 added a sanitized effective configuration snapshot to `benchmark.py`. Benchmark artifacts now begin with JSON context for model/API/dimensions, git state, Docker compose config/image identifiers, allowlisted env values, Redis INFO summary, environment baseline hash, and explicit redaction policy. The implementation preserves benchmark behavior, removes raw text output from the repeated request section, and corrects stale character-count labels. Verification produced `benchmark-results/fd-benchmark-m009-s01.txt` and proved the snapshot is present and safe.

## Verification

S01 passed compile, snapshot parser, Docker compose config, Go tests, full benchmark run, and artifact parser checks.

## Requirements Advanced

- R004 — Benchmark artifacts now record sanitized effective configuration context.
- R003 — Snapshot allowlist includes planned env-configurable cache/runtime fields.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Also corrected stale character-count labels in benchmark text labels because the new comparability artifact made the mismatch visible.

## Known Limitations

Docker image metadata is captured as compose image output lines plus hash rather than fully structured per-service objects. Benchmark still restarts API in Section 5 as before.

## Follow-ups

S02 should align cache namespace/retention env fields with the snapshot allowlist: Redis TTL/no-expire, model id/revision, cache/schema/tokenizer/chunking version, and dimension.

## Files Created/Modified

- `benchmark.py` — Adds sanitized effective config snapshot helpers, safe BENCHMARK_* overrides, corrected labels, and removes raw text printout.
- `benchmark-results/fd-benchmark-m009-s01.txt` — S01 benchmark artifact with snapshot and verification output.
