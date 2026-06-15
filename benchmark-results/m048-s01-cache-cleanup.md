# M048 S01 Cache Cleanup Evidence

Captured: 2026-06-15

## Scope

S01 covers GitHub issue #7 findings:

- #19 `LRUCache` is dead production code.
- #27 duplicate hash-truncate helpers in the cache package.
- #28 duplicated env integer parsers across cache/middleware/main.

Input issue artifact: `documents/issue-7-current-m048.md`.

## Pre-fix Evidence

Static proof:

```text
gsd_exec 12ffe3b3-84f7-4e6f-8f6a-e3a12f9eef57
PASS pre-fix issue #7 #19/#27/#28 cleanup debt present
```

Confirmed:

- `api/cache/lru.go` existed.
- Only non-self `LRUCache` reference was `api/fd_v2_cache_integration_test.go`.
- Both `shortCacheKeyHash` and `shortNamespaceHash` existed.
- Env integer parsers existed in `lru.go`, `ratelimit.go`, and `main.go`.

## Fix

- Deleted dead `api/cache/lru.go`, `api/cache/lru_test.go`, and `api/cache/lru_rapid_test.go`.
- Replaced fd_v2 cache integration's LRU scaffold with a LocalCache-backed test adapter that preserves cache HIT/MISS metrics assertions.
- Added `api/cache/hash.go` with canonical package-local `shortHash`.
- Updated TieredCache and Redis namespace code to use `shortHash`.
- Added `api/internal/envutil` with `Int` and `PositiveInt`.
- Updated `main.go` and `middleware/ratelimit.go` to use `envutil` and removed duplicate env parser functions.

## Green Evidence

Commands:

```bash
cd api && go test ./cache
cd api && go test ./...
```

Results:

```text
go test ./cache: 36 passed in 1 package
go test ./...: 282 passed in 10 packages
```

Post-cleanup static proof:

```text
gsd_exec 1453b735-d079-4ce7-9282-08805c13a318
PASS M048 S01 removed LRU and duplicate hash/env helpers
```

## Requirement Outcome

- R037 validated for issue #7 findings #19, #27, and #28.

## Residual Issue #7 Findings

Deferred to downstream M048 slices:

- #26 ONNX-only RuntimeHealth fields.
- #29 duplicate Embedder/WarmupModel interfaces.
- #30 lifecycle default singleton.
- #24 malformed validation message.
- #31 OpenAPI helper silent key drop.
