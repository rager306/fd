# S02: Cache namespace and retention — UAT

**Milestone:** M009-zjrq6j
**Written:** 2026-05-19T17:50:36.089Z

# UAT: S02 Cache namespace and retention

## Evidence

- Go tests passed: 60 tests in 4 packages.
- Pinned GolangCI-Lint passed: 0 issues.
- API image rebuilt and `/health` returned ok.
- Benchmark artifact `benchmark-results/fd-benchmark-m009-s02.txt` includes cache env values.
- Redis key after model-aware run used hashed namespace: `embed:cache:v2:m...:...:d1024`.
- Redis TTL for `REDIS_CACHE_TTL=168h` was positive (`604775` seconds at check time).
- No-expire mode produced `TTL == -1` and API was restored healthy afterward.

## Acceptance

- Defaults preserve legacy `v2` namespace and 24h TTL.
- Invalid TTL and TTL/no-expire conflict are rejected by tests.
- Model-aware namespace fields are opt-in and hashed in keys.
- Compose propagates env settings into the API container.

