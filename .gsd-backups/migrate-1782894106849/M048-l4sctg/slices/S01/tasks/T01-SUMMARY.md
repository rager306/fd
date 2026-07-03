---
id: T01
parent: S01
milestone: M048-l4sctg
key_files:
  - documents/issue-7-current-m048.md
  - api/cache/lru.go
  - api/cache/tiered.go
  - api/cache/redis.go
  - api/main.go
  - api/middleware/ratelimit.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T10:51:06.770Z
blocker_discovered: false
---

# T01: Confirmed issue #7 cache cleanup debt exists before S01 fixes.

**Confirmed issue #7 cache cleanup debt exists before S01 fixes.**

## What Happened

Ran a static pre-fix check proving `api/cache/lru.go` exists, has only a single non-self reference in `fd_v2_cache_integration_test.go`, duplicate short hash helpers exist in `tiered.go` and `redis.go`, and env integer parsers exist in `lru.go`, `ratelimit.go`, and `main.go`.

## Verification

Static proof `12ffe3b3-84f7-4e6f-8f6a-e3a12f9eef57` passed for pre-fix issue #7 findings #19/#27/#28.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_exec 12ffe3b3-84f7-4e6f-8f6a-e3a12f9eef57` | 0 | ✅ pass | 85ms |

## Deviations

GitNexus did not resolve package-local cache cleanup symbols; used direct static scan evidence instead.

## Known Issues

S01 remains to remove/unify the confirmed symbols.

## Files Created/Modified

- `documents/issue-7-current-m048.md`
- `api/cache/lru.go`
- `api/cache/tiered.go`
- `api/cache/redis.go`
- `api/main.go`
- `api/middleware/ratelimit.go`
