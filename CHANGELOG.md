# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] — 2026-05-14

### Added

- **L1 local cache** (`cache/local.go`) — sync.Map with per-entry TTL eviction, ~50ns access
- **Binary storage** (`cache/redis.go`) — `[dim:uint16][float32*dim]` encoding, 4 KB vs ~8 KB JSON for 1024d
- **Pool timeouts** — `DialTimeout: 5s`, `ReadTimeout: 3s`, `WriteTimeout: 3s`, `PoolTimeout: 4s`, `MinIdleConns: 10`
- **Two-tier cache** (`cache/tiered.go`) — L1 → L2 → loader with singleflight stampede prevention

### Changed

- **Storage format**: JSON HASH (`HSET/HGET`) → binary flat (`SET/GET`), 2x smaller
- **Cache architecture**: Redis-only → two-tier (sync.Map + Redis + singleflight)
- **Redis key format**: `embed:cache:text:<hash>` → `embed:cache:v2:<hash>:d<dim>` (dimension-aware)
- **Handler interface**: `cache.Get`/`Set` → `cache.GetOrLoad` (cache-aside pattern)

### Fixed

- Cold latency: ~70ms → ~19ms (TEI warmup + binary payload)
- 512d warm latency: ~9ms → ~2.6ms (was inflated by JSON overhead in Redis-only mode)
- Missing connection pool timeouts causing potential hanging connections

### Performance

| Metric | Before | After |
|--------|--------|-------|
| Warm latency (1024d) | 2.6 ms | 2.6 ms |
| Warm latency (512d) | 9.0 ms | 2.6 ms |
| Cold latency | ~70 ms | ~19 ms |
| Storage per embedding | ~8 KB | 4 KB |
| L1 cache | none | sync.Map (~50ns) |
| singleflight | none | stampede prevention |

## [0.1.0] — 2026-05-14

### Added

- Initial embedding service with TEI (deepvk/USER-bge-m3) + Redis Stack + Go API
- OpenAI-compatible `/v1/embeddings` endpoint
- `/embeddings/batch` endpoint for FalkorDB
- Matryoshka dimensions: 1024d (nodes), 512d (edges)
- Docker Compose setup with healthchecks
- Unit and integration tests

[0.2.0]: https://github.com/rager306/fd/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/rager306/fd/releases/tag/v0.1.0
