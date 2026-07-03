# S03: Redis persistence hardening — UAT

**Milestone:** M009-zjrq6j
**Written:** 2026-05-19T17:59:30.661Z

# UAT: S03 Redis persistence hardening

## Evidence

- `docker compose config` passed.
- Go tests passed: 60 tests in 4 packages.
- Pinned GolangCI-Lint passed: 0 issues.
- Redis CONFIG verification under explicit env:
  - `maxmemory = 134217728`
  - `maxmemory-policy = allkeys-lfu`
  - `save = 300 1`
  - `appendonly = no`
- Cached key survived Redis restart after `BGSAVE`.
- API restart after Redis restart succeeded and same embedding request returned ok.
- `benchmark-results/fd-benchmark-m009-s03.txt` includes Redis CONFIG snapshot fields.

## Acceptance

- Redis persistence is documented and configurable.
- Redis localhost binding remains in local override.
- Benchmark snapshot records Redis maxmemory/policy/persistence fields.
- Restart reuse was verified without relying on L1 cache.

