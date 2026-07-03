---
id: S02
parent: M009-zjrq6j
milestone: M009-zjrq6j
provides:
  - Env-configured Redis cache namespace and retention.
  - Validated Compose propagation for cache env settings.
  - S02 benchmark artifact with effective config snapshot.
requires:
  []
affects:
  []
key_files:
  - api/cache/redis.go
  - api/cache/redis_test.go
  - api/main.go
  - api/main_test.go
  - benchmark.py
  - docker-compose.override.yaml
  - benchmark-results/fd-benchmark-m009-s02.txt
key_decisions:
  - Preserve legacy `v2` keys and 24h TTL by default to avoid surprising cache invalidation.
  - Use opt-in hashed namespace fields for model/revision/tokenizer/chunking values.
  - Reject contradictory `REDIS_CACHE_TTL` plus `REDIS_CACHE_NO_EXPIRE=true` at startup.
  - Compose override must pass cache env vars into the API container; host benchmark env alone is not effective runtime config.
patterns_established:
  - Preserve defaults while making correctness-affecting cache namespace opt-in.
  - Hash correctness namespace values inside Redis keys; record clear values only in sanitized benchmark/config snapshots.
  - Validate cache config at startup instead of silently falling back.
observability_surfaces:
  - Benchmark snapshot now records MODEL_ID, REDIS_POOL_SIZE, REDIS_CACHE_TTL, and EMBEDDING_* fields when set. Redis key namespace can be correlated to the snapshot without raw model strings in keys.
drill_down_paths:
  - .gsd/milestones/M009-zjrq6j/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M009-zjrq6j/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M009-zjrq6j/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T17:50:36.088Z
blocker_discovered: false
---

# S02: Cache namespace and retention

**S02 made Redis cache retention and namespace env-configurable, tested, and visible in benchmark snapshots.**

## What Happened

S02 implemented env-configurable cache namespace and retention. Redis cache now has validated options, opt-in model-aware namespace hashing, configurable TTL, no-expire mode, and strict startup validation. The API startup path now fails fast on invalid cache config. Compose override propagates the new env knobs into the container. Tests cover defaults, validation, key namespace differences, and TTL/no-expire behavior. Runtime verification rebuilt the API, confirmed health, ran a benchmark with 168h/model-aware config, verified namespaced Redis keys and TTL, tested no-expire live, and restored the API to default env.

## Verification

S02 passed tests, lint, Docker compose config, API rebuild/health, benchmark run, artifact parser, Redis TTL/no-expire checks, and GitNexus detect_changes.

## Requirements Advanced

- R003 — Cache/runtime tuning parameters are env-configurable and validated.
- R004 — Benchmark artifacts record effective cache namespace and retention values.
- R002 — Long-lived cache retention modes now exist through TTL/no-expire settings.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Strict Redis pool validation exposed and fixed a pre-existing `getEnvInt` bug where an unset env var returned 0 instead of the default. S02 also updated Compose override to actually pass the new env knobs into the API container, because benchmark host env alone would not be an effective runtime config.

## Known Limitations

Default runtime remains legacy namespace unless model-aware env fields are explicitly set. Redis key namespace uses hashes, so operators need benchmark/config snapshot to map a key namespace back to model settings.

## Follow-ups

S03 should document Redis persistence/maxmemory policy using the new env/snapshot foundation. S04 can now assume benchmark artifacts record cache namespace and retention config.

## Files Created/Modified

- `api/cache/redis.go` — Adds Redis cache options, env parsing, opt-in namespace hashing, TTL/no-expire behavior, and key helper.
- `api/cache/redis_test.go` — Adds tests for default key preservation, namespace hashing, env parsing, TTL/no-expire, and invalid config.
- `api/main.go` — Wires validated Redis cache options into API startup and fixes getEnvInt default behavior.
- `api/main_test.go` — Adds tests for getEnvInt default, valid, and invalid behavior.
- `benchmark.py` — Adds new cache/runtime env fields to benchmark snapshot allowlist.
- `docker-compose.override.yaml` — Propagates cache namespace and retention env vars into the API container with safe defaults.
- `benchmark-results/fd-benchmark-m009-s02.txt` — Benchmark artifact proving S02 env/config snapshot and runtime behavior.
