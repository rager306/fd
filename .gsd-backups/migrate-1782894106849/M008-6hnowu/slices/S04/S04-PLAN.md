# S04: S04

**Goal:** Research current Redis practices that can improve throughput and latency for cached embedding workloads without changing the embedding model.
**Demo:** After this, Redis/cache throughput optimization opportunities are ranked for cached embeddings.

## Must-Haves

- Current Redis best practices are checked from Redis docs or credible sources.
- Candidate improvements cover pipelining/MGET, connection pool tuning, binary payload layout, compression tradeoffs, eviction/memory policy, client-side caching, clustering/IO threading where relevant, and observability.
- Each candidate is mapped to fd feasibility and benchmark method.
- Risky or irrelevant practices are explicitly excluded.

## Proof Level

- This slice proves: source research plus architecture assessment

## Integration Closure

Separates cache-layer improvements from inference runtime changes and maps them to benchmarkable fd changes.

## Verification

- Identifies Redis/API metrics needed to distinguish Redis bottlenecks from serialization, networking, and handler overhead.

## Tasks

- [x] **T01: Added Redis throughput research: fd should first benchmark MGET/pipelined batch cache hits and pool/round-trip metrics, not generic Redis tuning.** `est:small`
  Research Redis performance fundamentals relevant to fd cached embeddings: pipelining/MGET, connection pooling, latency diagnosis, network round trips, and client behavior. Use current Redis docs or credible sources.
  - Verify: Sources read and candidate practices recorded with relevance to fd.

- [x] **T02: Assess Redis embedding data layout and benchmarked env retention policy** `est:medium`
  Assess Redis data layout and retention policy for embedding vectors: binary value layout, float32 bytes, metadata/key design, long configurable TTL vs no-expire research mode, explicit invalidation, model/version-aware keys, compression tradeoffs, eviction policy, persistence settings, and Redis 7/8 features where relevant. Identify which parameters should become env vars, whether they affect cache correctness or only runtime performance, and which must be recorded in benchmark config snapshots.
  - Verify: Candidate data layout, retention/policy changes, env knobs, and benchmark-recorded fields are ranked with benchmark method and risks.

- [x] **T03: Research advanced Redis and benchmark-recorded deployment options** `est:medium`
  Research advanced Redis options for high-throughput and long-lived cached embedding workloads: client-side caching/tracking, Lua/functions if useful, clustering/sharding, I/O threading, persistence/RDB/AOF implications, Dragonfly/Valkey/Redis Stack alternatives, and when each is inappropriate for fd. Include which deployment/runtime knobs should be configurable via env and recorded in benchmark artifacts.
  - Verify: Advanced options are classified as candidate, future, or not recommended, with env-tunable and benchmark-recorded scope identified.

- [x] **T04: Recommend Redis reusable vector cache and comparable benchmark path** `est:small`
  Produce a Redis optimization recommendation for fd cached embeddings with proposed benchmark additions, success metrics, retention policy, repeated-chunk reuse test, env var configuration surface, sanitized benchmark configuration snapshot requirements, and stop criteria.
  - Files: `.gsd/milestones/M008-6hnowu/slices/S04/S04-RESEARCH.md`
  - Verify: Research artifact includes ranked options, benchmark plan, retention policy, env vars, sanitized config snapshot, and exclusions.

## Files Likely Touched

- .gsd/milestones/M008-6hnowu/slices/S04/S04-RESEARCH.md
