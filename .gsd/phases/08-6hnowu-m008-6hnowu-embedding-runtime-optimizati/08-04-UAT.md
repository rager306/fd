# S04: Research Redis cached embedding throughput — UAT

**Milestone:** M008-6hnowu
**Written:** 2026-05-19T17:10:15.219Z

# UAT: S04 Research Redis cached embedding throughput

## Evidence

- T01 mapped current cache round trips and batch-hit weakness.
- T02 assessed binary layout and long-retention policy.
- T03 ranked advanced Redis/deployment options.
- T04 produced recommendation artifact.

## Acceptance

- Redis L2 should be treated as a long-lived model-aware dense-vector cache.
- Next implementation should add retention/env/config snapshot and batch-hit benchmarks before optimizing code.
- Advanced Redis infrastructure remains deferred until measured bottleneck evidence exists.

