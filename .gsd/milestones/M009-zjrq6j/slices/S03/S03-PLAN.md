# S03: Redis persistence hardening

**Goal:** Document and configure an RDB-first Redis persistence path plus maxmemory/policy recording for long-lived embedding cache.
**Demo:** After this, local Redis can survive restarts for research cache reuse with documented memory/eviction policy.

## Must-Haves

- Redis persistence approach is documented and configurable.
- Redis remains localhost-bound for host access.
- Redis maxmemory/policy are documented and recorded by benchmark snapshot.
- Restart reuse is verified when persistence is enabled.

## Proof Level

- This slice proves: docker compose config plus Redis restart reuse check

## Integration Closure

Deployment docs and Compose/runtime settings match benchmark snapshot fields.

## Verification

- Benchmark and docs record Redis persistence, maxmemory, policy, dbsize, hit/miss, and eviction state.

## Tasks

- [x] **T01: Design Redis persistence config** `est:small`
  Inspect current Redis compose command, README runtime docs, and benchmark Redis metadata collection. Design minimal RDB-first persistence config and Redis CONFIG snapshot fields.
  - Files: `docker-compose.yaml`, `docker-compose.override.yaml`, `README.md`, `benchmark.py`
  - Verify: Summary identifies Redis config knobs, docs section, and benchmark fields.

- [x] **T02: Implement Redis persistence visibility** `est:medium`
  Implement Redis persistence/maxmemory/policy configuration and benchmark snapshot visibility. Keep Redis localhost binding safe. Document RDB-first cache persistence, maxmemory policy, and how to enable/disable AOF.
  - Files: `docker-compose.yaml`, `docker-compose.override.yaml`, `README.md`, `benchmark.py`
  - Verify: `docker compose config` and snapshot parser show Redis maxmemory/policy/persistence fields.

- [x] **T03: Verify Redis persistence hardening** `est:medium`
  Verify Redis restart reuse with persistence enabled: create a cached embedding, force/save if needed, restart Redis, verify key remains, restart API to clear L1, and confirm request succeeds from Redis-backed cache path. Run tests/lint/config and GitNexus detect_changes.
  - Files: `docker-compose.yaml`, `docker-compose.override.yaml`, `README.md`, `benchmark.py`
  - Verify: `docker compose config`; Go tests; lint; Redis restart reuse script; benchmark snapshot parser; GitNexus detect_changes.

## Files Likely Touched

- docker-compose.yaml
- docker-compose.override.yaml
- README.md
- benchmark.py
