---
id: T04
parent: S04
milestone: M008-6hnowu
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S04/S04-RESEARCH.md
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:09:55.369Z
blocker_discovered: false
---

# T04: Recommended next Redis spike: model-aware long-lived cache, sanitized config snapshots, and measured MGET/pipeline batch-hit benchmark.

**Recommended next Redis spike: model-aware long-lived cache, sanitized config snapshots, and measured MGET/pipeline batch-hit benchmark.**

## What Happened

Recommended Redis reusable vector cache and comparable benchmark path. The next implementation spike should keep the dense binary payload, add model/version-aware cache namespaces, make retention env-configurable with TTL/no-expire modes, document/deploy RDB persistence for long-lived cache survival, and extend benchmark output with sanitized effective config snapshots. Throughput work should focus on batch cache-hit workloads where current per-input GET calls multiply round trips; benchmark before/after with L1 disabled/reset, Redis INFO stats, go-redis pool stats, p50/p95/p99, RPS, hit/miss/eviction counts, memory, and Docker/runtime identifiers. Only after this evidence should fd test MGET/pipelined batch cache access. Cluster, Redis Stack, Dragonfly/Valkey, and advanced server-side features are not first-line changes.

## Verification

Saved S04 research artifact and aligned recommendation to R002/R003/R004/D003/D004/D005.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Saved: .gsd/milestones/M008-6hnowu/slices/S04/S04-RESEARCH.md` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Requirements advanced: R002, R003, R004` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Decisions applied: D003, D004, D005` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

This is a recommendation/research output; implementation still needs impact analysis for `RedisCache`, `TieredCache`, `CreateBatchEmbeddings`, and benchmark changes before editing symbols.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S04/S04-RESEARCH.md`
