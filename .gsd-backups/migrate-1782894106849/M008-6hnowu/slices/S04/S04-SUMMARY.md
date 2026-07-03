---
id: S04
parent: M008-6hnowu
milestone: M008-6hnowu
provides:
  - Redis long-lived cache design.
  - Batch cache-hit benchmark path.
  - Sanitized config snapshot field list.
requires:
  []
affects:
  []
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S04/S04-RESEARCH.md
key_decisions:
  - Redis L2 is a long-lived reusable dense-vector cache for research/chunking, not only a transient 24h cache.
  - Keep binary dense payload first; do not add compression without memory/network evidence.
  - Use RDB persistence first for re-embeddable long-lived cache survival; consider AOF only if data-loss cost is too high.
patterns_established:
  - Model/version-aware cache keys before long retention.
  - Measure Redis round trips and pool pressure before optimizing cache code.
  - Prefer boring Redis settings before replacing Redis infrastructure.
observability_surfaces:
  - Benchmark config snapshots should include Redis INFO stats, pool stats, memory, eviction policy, persistence mode, cache namespace, and benchmark corpus metadata.
drill_down_paths:
  - .gsd/milestones/M008-6hnowu/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S04/tasks/T02-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S04/tasks/T03-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S04/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T17:10:15.219Z
blocker_discovered: false
---

# S04: Research Redis cached embedding throughput

**S04 produced the Redis vector-cache plan: long-lived model-aware cache plus comparable batch-hit benchmarks.**

## What Happened

S04 researched Redis cache throughput and long-lived reusable embedding cache design. It confirmed current dense binary layout is sound, identified fixed TTL and incomplete namespace as main correctness/retention gaps, ranked MGET/pipelining as the likely batch cache-hit throughput improvement, and defined benchmark snapshot fields for comparable future runs. Advanced Redis deployment options were intentionally deferred until measured bottlenecks justify them.

## Verification

All S04 tasks complete and S04 research artifact saved.

## Requirements Advanced

- R002 — Defined long-lived retention, persistence, and reuse strategy.
- R003 — Defined env/config knobs for cache namespace, TTL, Redis policy, and batch mode.
- R004 — Defined sanitized benchmark configuration snapshot fields.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S04 scope expanded after user clarified Redis should support long-lived research/chunking reuse and env-configurable benchmark comparability.

## Known Limitations

No cache implementation changed in this research slice. Current code still has fixed 24h Redis TTL and no automatic sanitized config snapshot.

## Follow-ups

Implementation spike should add env-configurable retention/namespace settings, sanitized benchmark config snapshots, and Redis batch-hit benchmark before changing cache code. Run GitNexus impact before editing RedisCache/TieredCache/CreateBatchEmbeddings.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S04/S04-RESEARCH.md` — Redis cache research and recommendation artifact.
