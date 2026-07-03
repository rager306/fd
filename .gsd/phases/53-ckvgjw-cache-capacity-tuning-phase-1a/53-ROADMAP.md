# M053-ckvgjw: Cache capacity tuning (Phase 1a)

**Vision:** Реализовать Phase 1a из Issue #9 — env-tunable L1 cache capacity и TTL. Сейчас `cache.NewLocalCache(10000, 30*time.Second)` hardcoded в main.go. После завершения: операторы могут увеличить L1 capacity (например до 1M для 4GB RAM) и TTL без recompiling fd-api. Default values сохраняют текущее поведение — backward compatible. Phase 1b (miss-coalescing) и Phase 1c (TEI thread-tuning) deferred до data-driven decision baseline и README guidance соответственно.

Не меняем: LocalCache API, TieredCache API, embed.Embedder interface, /health shape, /v1/embeddings behavior. Cache flush endpoint (M049) работает с новым config. Recovery contract (M051) ортогонален.

## Slices

- [x] **S01: L1 cache capacity + TTL env-config** `risk:low` `depends:[]`
  > After this: FD_CACHE_MAX_SIZE=100 FD_CACHE_LOCAL_TTL=60s docker compose up — после мониторинга /metrics (fd_cache_entries{tier=l1} ≈ FD_CACHE_MAX_SIZE) и поведения flush endpoint.

## Boundary Map

Not provided.
