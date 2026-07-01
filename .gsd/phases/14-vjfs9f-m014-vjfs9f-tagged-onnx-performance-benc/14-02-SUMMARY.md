---
id: S02
parent: M014-vjfs9f
milestone: M014-vjfs9f
provides:
  - TEI baseline metrics and artifact for tagged ONNX comparison.
requires:
  []
affects:
  - S03
  - S04
key_files:
  - benchmark-results/fd-benchmark-m014-tei-baseline.txt
key_decisions:
  - Use `benchmark-results/fd-benchmark-m014-tei-baseline.txt` as M014 TEI baseline.
  - Accept git dirty=true in active-slice benchmark snapshot as documented context, not a benchmark invalidation.
patterns_established:
  - Benchmark artifacts are verified for actual headings, required metadata, and raw-text hygiene before slice completion.
observability_surfaces:
  - TEI benchmark artifact snapshot v2.
  - Redis before-run stats and benchmark cache deltas.
  - Docker compose image/config hashes.
drill_down_paths:
  - .gsd/milestones/M014-vjfs9f/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M014-vjfs9f/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M014-vjfs9f/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T04:20:32.936Z
blocker_discovered: false
---

# S02: TEI baseline benchmark

**S02 captured and verified the fresh TEI control benchmark for M014.**

## What Happened

S02 verified the default TEI stack and ran a fresh TEI benchmark through `benchmark.py` snapshot v2. The artifact records Docker/Redis/git/environment metadata, runtime label `tei-default`, cold/warm latency, cache-hit behavior, concurrent throughput, Redis L2 persistence after API restart, cached batch behavior, and repeated chunk reuse. Artifact hygiene checks passed with no raw fixed-probe text leaks.

## Verification

Preflight, benchmark run, parser/leak checks, and GitNexus artifact-scope verification passed.

## Requirements Advanced

- benchmark-comparability — Produced a comparable TEI baseline with sanitized config snapshot.

## Requirements Validated

- tei-baseline-control — `benchmark-results/fd-benchmark-m014-tei-baseline.txt` includes snapshot_version 2 and runtime_label `tei-default`; hygiene checks passed.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

The artifact parser initially used stale expected heading names and was corrected to match actual benchmark output. The benchmark artifact records git dirty=true because GSD artifacts were open during the run.

## Known Limitations

Benchmark uses the current local Docker/KVM host; not a CI or production-host result. Snapshot records TEI as default runtime and no ONNX/native metadata.

## Follow-ups

S03 should run tagged ONNX benchmark with isolated Redis namespace and comparable snapshot fields. S04 should compare against TEI results: best cold 59.0ms, warm mean 2.25ms, max throughput ~750 req/s at 16 concurrent, batch L1 p95 4.16ms, batch L2 p95 5.51ms, chunk reuse warm p95 7.04ms.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m014-tei-baseline.txt` — Fresh TEI baseline benchmark artifact with snapshot v2 metadata.
