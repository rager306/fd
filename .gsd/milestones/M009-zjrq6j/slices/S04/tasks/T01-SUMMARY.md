---
id: T01
parent: S04
milestone: M009-zjrq6j
key_files:
  - benchmark.py
  - api/handlers/batch.go
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T18:02:19.430Z
blocker_discovered: false
---

# T01: Designed S04 as benchmark-only helpers and sections for L1, Redis L2, cached batch, repeated chunk reuse, and Redis INFO deltas.

**Designed S04 as benchmark-only helpers and sections for L1, Redis L2, cached batch, repeated chunk reuse, and Redis INFO deltas.**

## What Happened

Designed S04 benchmark sections. Add helpers: `call_batch_api(inputs)`, `redis_stats_snapshot()`, `redis_stats_delta(before, after)`, `print_redis_delta(label, delta)`, and a reusable latency summary printer. Preserve existing sections but add Redis INFO deltas around repeated request hot hits and Redis L2 restart checks. Add a new cached batch section using `/embeddings/batch` with `inputs`, `dimensions`, and `encoding_format=base64`, priming deterministic labels and measuring repeated cached batch calls. Add a repeated chunk reuse section using a small set of repeated chunk labels across multiple batch requests to model research/chunking reuse. The artifact should print labels/counts only, never raw texts. Redis delta fields: keyspace_hits, keyspace_misses, evicted_keys, expired_keys, total_commands_processed if available.

## Verification

Read current benchmark and batch handler. GitNexus impact for benchmark `main` and `call_api` was LOW with only benchmark-file callers affected.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Read: benchmark.py` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read: api/handlers/batch.go` | -1 | unknown (coerced from string) | 0ms |
| 3 | `GitNexus impact: Function:benchmark.py:main LOW; call_api LOW` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

Existing benchmark Section 2 is effectively an L1 hot-hit workload after priming, but it does not label Redis INFO deltas. Existing Section 5 exercises Redis L2 after API restart, but also lacks hit/miss deltas. S04 will add explicit diagnostics rather than change service behavior.

## Files Created/Modified

- `benchmark.py`
- `api/handlers/batch.go`
