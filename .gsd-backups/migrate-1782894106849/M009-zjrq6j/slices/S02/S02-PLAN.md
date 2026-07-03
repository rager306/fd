# S02: Cache namespace and retention

**Goal:** Expose model-aware cache namespace and TTL/no-expire retention through safe env configuration while preserving defaults.
**Demo:** After this, Redis cache retention and namespace can be tuned safely for long-lived research reuse.

## Must-Haves

- Defaults preserve current behavior.
- Env validation rejects invalid or contradictory cache retention settings.
- Cache keys include explicit model/schema/version namespace fields.
- TTL/no-expire behavior is tested.

## Proof Level

- This slice proves: Go tests for defaults, env validation, namespace differences, TTL/no-expire behavior

## Integration Closure

Runtime cache keys and retention align with benchmark snapshot fields from S01.

## Verification

- Cache configuration becomes inspectable through benchmark snapshots and safe logs without raw text or secrets.

## Tasks

- [x] **T01: Design cache config surface** `est:small`
  Inspect current Redis/Tiered cache construction and key generation. Design a small cache config surface with defaults preserving current keys/TTL unless env vars are set, plus validation for TTL/no-expire conflicts.
  - Files: `api/cache/redis.go`, `api/cache/tiered.go`, `api/main.go`, `benchmark.py`
  - Verify: Summary names defaults, env vars, validation rules, and affected symbols.

- [x] **T02: Implement cache namespace and retention** `est:medium`
  Implement cache config parsing and Redis key namespace/retention behavior. Preserve current default TTL and key behavior as closely as possible, support explicit namespace fields, and support no-expire mode. Wire from main runtime env and keep benchmark snapshot allowlist aligned.
  - Files: `api/cache/redis.go`, `api/cache/redis_test.go`, `api/main.go`, `benchmark.py`
  - Verify: Targeted Go tests and Python snapshot parser pass.

- [x] **T03: Verify cache config surface** `est:medium`
  Run S02 verification: Go tests, pinned GolangCI-Lint, docker compose config, and benchmark snapshot parser confirming new env fields are included safely. Run GitNexus detect_changes before slice completion/commit.
  - Files: `api/cache/redis.go`, `api/main.go`, `benchmark.py`
  - Verify: `cd api && go test ./... -short`; pinned GolangCI-Lint; `docker compose config`; snapshot parser with env fields; GitNexus detect_changes.

## Files Likely Touched

- api/cache/redis.go
- api/cache/tiered.go
- api/main.go
- benchmark.py
- api/cache/redis_test.go
