---
id: T03
parent: S04
milestone: M008-6hnowu
key_files:
  - docker-compose.yaml
  - docker-compose.override.yaml
  - README.md
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:06:00.600Z
blocker_discovered: false
---

# T03: Ranked Redis options: MGET/pipeline, pool metrics, long retention, RDB, and maxmemory policy first; cluster/Dragonfly/Stack later only if bottleneck proves it.

**Ranked Redis options: MGET/pipeline, pool metrics, long retention, RDB, and maxmemory policy first; cluster/Dragonfly/Stack later only if bottleneck proves it.**

## What Happened

Researched advanced Redis and deployment options for reusable vector cache. Low-risk candidates are batch-oriented access (`MGET` or bounded pipelining), go-redis pool stats/latency observability, explicit Redis `maxmemory` plus `allkeys-lru` or `allkeys-lfu`, RDB persistence for research cache survival, and env-configured TTL/no-expire modes. Medium/future candidates include Unix socket or colocated network path tests if Docker bridge latency shows up, read replicas only for multi-node read scaling, and Redis server-assisted client-side caching only if multiple API replicas need invalidation beyond TTL. Infrastructure alternatives like Dragonfly/Valkey, Redis Cluster/sharding, Redis Stack/vector search, Lua/functions, and I/O threading should be excluded from the first spike: they either solve server-scale bottlenecks not yet proven, add operational complexity, or address vector search rather than key-value embedding cache. Deployment knobs to record in benchmarks include Redis image/version, bind mode, persistence config, maxmemory/policy, pool settings, batch mode, Redis INFO memory/stats, and Docker network/container limits.

## Verification

Read Redis docs and source summaries for pipelining, client-side caching, eviction, persistence, and current compose/README context from prior milestones.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Fetched: https://redis.io/docs/latest/develop/use/pipelining/` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Fetched: https://redis.io/docs/latest/develop/use/client-side-caching/` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Fetched: https://redis.io/docs/latest/develop/reference/eviction/` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Fetched: https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/` | -1 | unknown (coerced from string) | 0ms |
| 5 | `Read project files earlier: docker-compose.yaml, docker-compose.override.yaml, README.md` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

Advanced infrastructure options are likely premature until fd proves Redis is the bottleneck with per-layer timing, pool stats, and batch-hit benchmarks. Redis server-assisted client-side caching is mostly redundant with fd's existing L1 unless multi-replica invalidation becomes a real requirement.

## Files Created/Modified

- `docker-compose.yaml`
- `docker-compose.override.yaml`
- `README.md`
