# S04: Redis batch hit benchmark — UAT

**Milestone:** M009-zjrq6j
**Written:** 2026-05-19T18:13:31.806Z

# UAT: S04 Redis batch hit benchmark

## Evidence

- Go tests passed: 60 tests in 4 packages.
- Pinned GolangCI-Lint passed: 0 issues.
- `docker compose config` passed.
- Redis/API were force-recreated with default config before final artifact.
- `benchmark-results/fd-benchmark-m009-s04.txt` contains:
  - `redis_delta/l1_hot_repeated`
  - `redis_delta/l2_after_api_restart`
  - `redis_delta/batch_l1_hot`
  - `redis_delta/batch_l2_after_api_restart`
  - `redis_delta/repeated_chunk_reuse`
- Batch L2 after API restart: 16 Redis hits, 0 misses, p95 about 5.61ms.
- Warm repeated chunk reuse p95 about 4.22ms.

## Acceptance

- Benchmark isolates L1 hot hits and Redis L2 after API restart.
- Benchmark includes cached batch workloads and repeated chunk reuse.
- Redis diagnostic deltas are printed for each relevant section.
- Evidence is comparable with snapshot from S01-S03.

