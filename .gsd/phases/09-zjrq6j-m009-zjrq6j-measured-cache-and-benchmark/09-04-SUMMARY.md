---
id: S04
parent: M009-zjrq6j
milestone: M009-zjrq6j
provides:
  - Batch cache-hit benchmark sections.
  - Redis delta evidence for MGET/pipeline go/no-go.
  - S04 benchmark artifact.
requires:
  []
affects:
  []
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m009-s04.txt
key_decisions:
  - Separate first cold chunk round from warm repeated chunk reuse p95.
  - Use Redis INFO deltas to distinguish L1 hits from Redis L2 hits.
  - Do not treat host/container config mismatch as acceptable; force-recreate for comparable artifacts.
patterns_established:
  - Use Redis deltas to prove whether a workload hit L1 or L2.
  - Separate cold first-round cost from warm reuse cost.
  - Benchmark effective server config before interpreting performance.
observability_surfaces:
  - Redis INFO deltas per benchmark section, batch L1/L2 p95, chunk first-cold and warm-reuse p95, benchmark artifact with Redis CONFIG snapshot.
drill_down_paths:
  - .gsd/milestones/M009-zjrq6j/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M009-zjrq6j/slices/S04/tasks/T02-SUMMARY.md
  - .gsd/milestones/M009-zjrq6j/slices/S04/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T18:13:31.806Z
blocker_discovered: false
---

# S04: Redis batch hit benchmark

**S04 added and verified batch cache-hit benchmark evidence; MGET/pipeline is not justified by current local measurements.**

## What Happened

S04 extended the benchmark to explicitly measure cache-hit workloads. The artifact now reports Redis deltas for L1-hot repeated requests, Redis L2 after API restart, cached batch L1, cached batch L2 after API restart, and repeated chunk reuse. Results show the current batch cache-hit path is low milliseconds: 16-item Redis L2 batch p95 around 5.61ms with 16 Redis hits and zero misses; warm repeated chunk reuse p95 around 4.22ms after a cold first round. The main latency cost remains cold model calls, not Redis round trips at this scale.

## Verification

S04 passed tests, lint, compose config, benchmark run, parser checks, and GitNexus detect_changes.

## Requirements Advanced

- R004 — Benchmark artifacts now include section-level Redis diagnostics for comparable cache experiments.
- R002 — Repeated chunk reuse behavior is measured against long-lived cache assumptions.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S04 reran the benchmark after force-recreating Redis/API with default Compose settings because the first artifact inherited S03's explicit Redis test config. This preserved benchmark comparability.

## Known Limitations

S04 synthetic batch sizes are small and local. Larger production chunk batches may still need MGET/pipeline analysis later. Current evidence is enough to avoid speculative optimization in this milestone.

## Follow-ups

S05 should be skipped or deferred: S04 did not show Redis round-trip pressure large enough to justify MGET/pipeline now. If pursued later, use S04 artifact as baseline and target larger batch sizes or multi-request workloads.

## Files Created/Modified

- `benchmark.py` — Adds batch endpoint benchmark helper, latency summaries, Redis INFO deltas, cached batch sections, and repeated chunk reuse section.
- `benchmark-results/fd-benchmark-m009-s04.txt` — Benchmark artifact with S04 Redis delta evidence.
