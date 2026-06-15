# M049 S01 Cache Invalidation Evidence

Captured: 2026-06-15

## Scope

S01 implements GitHub issue #8 AN-A: authenticated cache invalidation for fd's embedding cache.

## Implemented

- `LocalCache.Flush(ctx)` clears all L1 entries.
- `LocalCache.Size()` reports current L1 entries.
- `RedisCache.Delete(ctx, input, dim)` deletes a single namespace/dimension-scoped Redis embedding entry.
- `RedisCache.FlushNamespace(ctx)` scans and deletes only keys matching the configured fd cache namespace.
- `TieredCache.Delete(ctx, input, dim)` removes one embedding entry from L1 and L2.
- `TieredCache.Flush(ctx)` clears L1 and fd's L2 namespace.
- `TieredCache.LocalSize()` exposes L1 occupancy for later metrics.
- `POST /v1/cache/flush` flushes the fd cache namespace.
- `POST /v1/cache/delete` deletes cache entries for `input` plus `dimensions`, defaulting dimensions to 1024.

## Safety Notes

- Redis invalidation never calls `FLUSHDB`.
- Redis flush is namespace-scoped via `prefix + namespace + ":*"`.
- Routes are registered after the existing global `APIKeyAuthFromEnv` middleware, so they are protected by the same bearer auth as other protected operator endpoints.
- The HTTP delete route uses input text plus dimensions rather than `:keyHash`; current cache keys are derived from input+dimension and short hashes are not reversible identifiers.

## Verification

Red test command:

```bash
cd api && go test ./cache ./handlers
```

Expected red result before implementation:

```text
missing LocalCache.Flush/Size, RedisCache.namespacePattern, TieredCache.Delete/Flush/LocalSize, CacheHandler/NewCacheHandler
```

Green commands:

```bash
cd api && go test ./cache ./handlers
cd api && go test ./...
```

Results:

```text
go test ./cache ./handlers: 127 passed in 2 packages
go test ./...: 293 passed in 10 packages
```

Static proof:

```text
gsd_exec 3670b28f-8bce-433e-8306-987102db98cb
PASS M049 S01 cache invalidation invariants
```

## Requirement Outcome

- R040 advanced: cache invalidation primitives and HTTP routes are implemented and tested.
- Runtime proof for authenticated live HIT->flush->MISS is deferred to M049 S03 after S02 is complete and the container is rebuilt once for the aggregate milestone.
