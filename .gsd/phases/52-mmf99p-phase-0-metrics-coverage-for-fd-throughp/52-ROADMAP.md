# M052-mmf99p: Phase 0 metrics coverage for fd throughput

**Vision:** Реализовать observability foundation из Issue #9 Phase 0 (cache hit-rate by tier, per-stage latency, TEI saturation, batch-fill ratio, cache occupancy) — инструментация которая нужна для data-driven tuning всех последующих фаз оптимизации (Phase 1 cache capacity/coalescing/thread-tuning, Phase 2 bulk queue). Цель: через `/metrics` видеть hit-rate L1/L2/miss, p95/p99 TEI inference latency, batch-fill относительно 32-cap, размер L1 и L2 namespaces, in-flight TEI calls и errors. Никаких изменений в runtime behavior — additive только. FD-сам contract, TEI backends, lifecycle state и embed.Embedder interface не трогаем.

## Slices

- [x] **S01: Phase 0 observability instrumentation** `risk:low` `depends:[]`
  > After this: /metrics endpoint экспонирует fd_cache_hits_total{tier=l1|l2|miss}, fd_tei_request_duration_seconds (p95/p99 видимы), fd_tei_batch_fill_ratio, fd_cache_entries{tier=l1|l2}, fd_cache_memory_bytes{tier} — без regression в существующих метриках.

## Boundary Map

Not provided.
