# S04 Research: Redis Cached Embedding Throughput and Reusable Vector Cache

## Scope

Research Redis/cache practices that can improve throughput and reuse for cached embeddings without changing the embedding model.

Key user constraints:

- Research/chunking workflows may reuse the same vectors multiple times.
- Redis L2 should be a sufficiently long-lived reusable embedding cache, not only a transient 24h hot cache.
- Many tuning parameters should be env-configurable.
- Benchmark artifacts must record sanitized effective config snapshots for comparability.

## Current fd cache behavior

Relevant files:

- `api/cache/redis.go`
- `api/cache/tiered.go`
- `api/handlers/embeddings.go`
- `api/handlers/batch.go`
- `benchmark.py`

Current behavior:

- L1 local cache stores binary embeddings.
- Redis L2 stores binary embeddings.
- Binary dense format: `[dim:uint16][float32*dim]`.
- 1024d vector is ~4098 bytes.
- Redis key currently includes schema-ish prefix and dimension suffix.
- Redis TTL is fixed at 24h.
- Batch handlers loop input-by-input through `GetOrLoad`, producing N potential Redis lookups.

## Findings

### 1. Keep dense binary layout first

The current binary payload is already a good default:

- compact compared with JSON;
- deterministic;
- easy to decode;
- avoids raw text in Redis values.

Do not compress first. Compression should only be tested if Redis memory or network bandwidth is proven bottleneck because compression adds CPU to every hit.

### 2. Add model/version-aware cache namespace

For long-lived research cache, key correctness matters more than raw hit rate. Required namespace fields:

- `EMBEDDING_MODEL_ID`
- `EMBEDDING_MODEL_REVISION`
- tokenizer id/hash/version
- embedding/cache schema version
- normalization/pooling version
- dimensions
- chunking version/salt when caching chunk embeddings

This prevents silent reuse of stale vectors after model/tokenizer/chunking changes.

### 3. Add configurable retention modes

Current fixed `24 * time.Hour` TTL is too short for repeated research/chunking workflows.

Recommended env surface:

- `REDIS_CACHE_TTL`: duration such as `24h`, `168h`, `720h`.
- `REDIS_CACHE_NO_EXPIRE`: boolean, mutually exclusive with TTL if enabled.
- `REDIS_MAXMEMORY`: deployment/server config, recorded by benchmark.
- `REDIS_MAXMEMORY_POLICY`: `allkeys-lru` default or `allkeys-lfu` for repeated chunk access.

Avoid `volatile-*` eviction policies when no-expire mode is enabled because no-expire keys may become non-evictable.

### 4. Persistence recommendation

Redis docs classify persistence options as RDB, AOF, no persistence, or RDB+AOF.

For fd embedding cache:

- RDB first: compact, faster restart, acceptable if losing a few minutes of re-embeddable cache writes is tolerable.
- AOF later: only if re-embedding loss is too expensive and write overhead is acceptable.
- no persistence only for ephemeral benchmark runs.

### 5. Throughput optimization path

Most likely near-term bottleneck in cache-hit batch workloads is per-input Redis round-trip.

Recommended implementation spike order:

1. Add benchmark sections that isolate L1 hit, Redis L2 hit, Redis miss/model call, and Redis restart persistence reuse.
2. Add batch cache-hit test with repeated deterministic inputs and L1 disabled or reset.
3. Implement bounded `MGET` or pipelined `GET` for batch dense cache hits.
4. Add go-redis pool stats and Redis INFO snapshots to benchmark artifacts.
5. Compare before/after at p50/p95/p99, RPS, Redis ops/sec, pool wait count, hit/miss counts, and memory.

### 6. Advanced options not first

Defer until bottleneck evidence exists:

- Redis Cluster/sharding: only if memory/throughput exceeds one node.
- Dragonfly/Valkey: only if Redis server CPU/latency is proven bottleneck.
- Redis Stack/vector search: fd currently needs key-value embedding cache, not vector database search.
- Redis server-assisted client-side caching: possibly useful with many API replicas, but current L1 already covers single-process hot hits.
- Lua/functions: not needed for simple get/set binary vectors.
- Unix socket/host networking: test only if Docker bridge latency appears material.

## Benchmark config snapshot fields

Every benchmark artifact should include sanitized effective configuration:

- git commit/branch/dirty flag;
- Docker image IDs and compose file hash;
- API env: `PORT`, `LOG_LEVEL`, cache TTL/no-expire, L1 size/TTL, batch cache mode, pipeline size;
- embedding namespace fields: model id/revision, tokenizer hash, schema version, dimensions;
- Redis config: version/image, bind mode, maxmemory, policy, persistence mode, dbsize, INFO memory/stats hit/miss/evictions;
- benchmark corpus name/version and input count;
- hardware baseline artifact reference, e.g. `benchmark-results/fd-environment-inxi-m008.txt`.

Secrets and raw texts must not be emitted.

## Recommendation

Treat Redis L2 as a long-lived, model-aware reusable dense-vector cache. First implementation spike should add env-configured retention/namespace/config snapshots and a batch Redis-hit benchmark; then test MGET/pipelining. Do not change model/runtime until Redis cache comparability and quality gates are in place.
