# S01: Cache cleanup consolidation — UAT

**Milestone:** M048-l4sctg
**Written:** 2026-06-15T11:05:02.471Z

# S01: Cache cleanup consolidation — UAT

**Milestone:** M048-l4sctg
**Written:** 2026-06-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S01 removes backend dead code and duplicate helpers. The observable outcome is source shape plus tests; no browser surface is involved.

## Preconditions

- `benchmark-results/m048-s01-cache-cleanup.md` exists.

## Smoke Test

Verify LRU removal, hash helper unification, envutil usage, and S01 evidence.

## Test Cases

### 1. LRU source removed

1. Inspect `api/cache/`.
2. **Expected:** `lru.go`, `lru_test.go`, and `lru_rapid_test.go` are absent.

### 2. Hash helper unified

1. Inspect `api/cache/hash.go`, `api/cache/tiered.go`, and `api/cache/redis.go`.
2. **Expected:** `shortHash` exists; `shortCacheKeyHash` and `shortNamespaceHash` are absent.

### 3. Active integer configuration parsing uses shared helper

1. Inspect `api/main.go`, `api/middleware/ratelimit.go`, and `api/internal/envutil/int.go`.
2. **Expected:** active parsing uses envutil helpers and duplicate parser functions are absent.

### 4. Evidence artifact complete

1. Inspect `benchmark-results/m048-s01-cache-cleanup.md`.
2. **Expected:** artifact covers #19, #27, #28, tests, and R037 validation.

## Edge Cases

- fd_v2 cache integration still proves second request is a cache HIT and metrics include hit counter.
- Envutil preserves zero-allowed semantics for main and positive-only semantics for rate limits.

## Failure Signals

- LRU symbols return.
- Duplicate hash helper names return.
- `go test ./...` fails.

## Requirements Proved By This UAT

- R037: dead/duplicate cache helper and env parsing cleanup.

## Not Proven By This UAT

- Runtime health/interface cleanup and API polish; these are S02/S03.

## Notes for Tester

UAT evidence: `2ae5e91d-6c8a-48f7-9e82-505921af6680`, `e7475039-5ae6-4261-b8cb-b3e48ad50841`, `1453b735-d079-4ce7-9282-08805c13a318`, `c26cc783-387c-41dd-8dbe-13c521b29e34`.
