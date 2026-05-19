---
id: M004-9886ht
title: "Measured performance optimization"
status: complete
completed_at: 2026-05-19T10:34:55.582Z
key_decisions:
  - Use debug-level cache path events and retain warn-level Redis degradation events.
  - Remove handler success INFO logs instead of sampling them because cache layer now owns cache-path diagnostics.
  - Keep Redis L2 restart diagnostic in benchmark.py because the benchmark already assumes local Docker/Redis control.
key_files:
  - benchmark.py
  - api/main.go
  - api/cache/tiered.go
  - api/cache/tiered_cache_test.go
  - api/handlers/embeddings.go
  - api/handlers/batch.go
  - api/handlers/embeddings_integration_test.go
  - benchmark-results/fd-benchmark-m004-final.txt
  - benchmark-results/fd-benchmark-m004-s01.txt
  - benchmark-results/fd-benchmark-m004-s03.txt
lessons_learned:
  - Benchmark summaries should be derived from stored measurement rows, not loop-local final values.
  - Default INFO logs should not scale linearly with high-throughput success paths.
  - Cache path debug logs must avoid raw input text; short hashes plus dimension are enough for correlation.
---

# M004-9886ht: Measured performance optimization

**Fixed benchmark throughput reporting, added cache-path observability, removed noisy success INFO logs, and verified Redis L2 restart behavior.**

## What Happened

M004 turned the M003 runtime baseline into measured performance and observability improvements. The benchmark max-throughput summary now selects the true measured row, eliminating the M003 contradiction. API runtime logs are quieter by default: success INFO logs from single and batch embedding handlers were removed, while invalid-request warnings and embedding errors remain. Cache-path observability moved to TieredCache with debug events for L1 hits, L2 hits, misses, dimension mismatches, and singleflight sharing, plus warn events for Redis get/set degradation without logging raw input text. The benchmark now includes a Redis L2 after API restart diagnostic; evidence showed Redis L2 served the primed request in 3.10ms after API restart in S03, and the final benchmark stayed self-consistent. Final verification passed across Compose config, Go tests, uv Python 3.13 benchmark, Compose health, and GitNexus change detection.

## Success Criteria Results

- Benchmark summary bug from M003 fixed: pass.
- Cache/runtime observability improved: pass.
- Default log noise reduced: pass.
- Redis L2 benchmark diagnostic added: pass.
- uv Python 3.13 evidence generated: pass.
- Final verification gates passed: pass.

## Definition of Done Results

- Benchmark summary reports actual max throughput: met.
- Cache/handler observability distinguishes cache paths without default success spam: met.
- Benchmark measures Redis-persisted cache after API restart under uv Python 3.13: met.
- Tests and runtime verification passed: met.
- GitNexus change detection low risk: met.

## Requirement Outcomes

No formal requirement IDs were changed. Operational/performance capability improved: benchmark evidence is more trustworthy, logs are more useful under load, and cache behavior is inspectable without leaking raw input text.

## Deviations

GitNexus did not directly resolve `TieredCache.GetOrLoad` as an indexed symbol; cache changes were covered by direct tests and runtime verification. Benchmark now restarts API during diagnostics, which is useful locally but intrusive in shared environments.

## Follow-ups

Potential next work: add rate limiting/counters for Redis degradation warnings if production Redis outages create noisy warn logs; consider fixing the existing `Cache speedup ... (median text)` wording because it currently reports max speedup rather than median-text speedup.
