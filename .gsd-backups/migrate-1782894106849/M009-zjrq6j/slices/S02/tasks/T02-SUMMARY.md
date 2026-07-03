---
id: T02
parent: S02
milestone: M009-zjrq6j
key_files:
  - api/cache/redis.go
  - api/cache/redis_test.go
  - api/main.go
  - benchmark.py
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T17:40:42.529Z
blocker_discovered: false
---

# T02: Implemented opt-in model-aware Redis namespace and env-configured TTL/no-expire retention while preserving default `v2` and 24h behavior.

**Implemented opt-in model-aware Redis namespace and env-configured TTL/no-expire retention while preserving default `v2` and 24h behavior.**

## What Happened

Implemented cache namespace and retention configuration. `RedisCacheOptions` and `RedisCacheNamespace` now control prefix, pool size, TTL/no-expire mode, and correctness-affecting namespace fields. Defaults preserve current behavior: `v2` namespace and 24h TTL. `RedisCacheOptionsFromEnv` supports `REDIS_CACHE_TTL`, `REDIS_CACHE_NO_EXPIRE`, `EMBEDDING_CACHE_VERSION`, `EMBEDDING_MODEL_ID`, `EMBEDDING_MODEL_REVISION`, `EMBEDDING_TOKENIZER_VERSION`, and `EMBEDDING_CHUNKING_VERSION`, rejecting invalid TTLs and TTL/no-expire conflicts. Namespace values are hashed before entering Redis keys, so raw model/tokenizer values are not exposed through key scans. `main.go` now validates cache config at startup and fails fast on invalid settings. `benchmark.py` snapshot allowlist now includes `MODEL_ID` and `REDIS_POOL_SIZE` so runtime/cache context remains comparable.

## Verification

Targeted Go cache tests and benchmark snapshot env-field check passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./cache -run 'TestRedisCache|TestHashText' -count=1` | 0 | ✅ pass: 10 tests in cache package | 8300ms |
| 2 | `REDIS_CACHE_TTL=168h EMBEDDING_MODEL_ID=deepvk/USER-bge-m3 REDIS_POOL_SIZE=50 MODEL_ID=deepvk/USER-bge-m3 uv run --python 3.13 --with requests --with redis python - <<'PY' ...` | 0 | ✅ pass: snapshot env fields check passed | 6700ms |

## Deviations

GitNexus symbol impact for Go cache symbols remains unavailable due GitNexus lookup issue, so direct references were enumerated with `rg`; `gitnexus_detect_changes` will be run before commit.

## Known Issues

Default runtime does not set `EMBEDDING_MODEL_ID`, so default keys intentionally remain legacy `v2` to avoid invalidating existing Redis caches. Model-aware namespace is opt-in through env vars.

## Files Created/Modified

- `api/cache/redis.go`
- `api/cache/redis_test.go`
- `api/main.go`
- `benchmark.py`
