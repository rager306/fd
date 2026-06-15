---
id: T03
parent: S01
milestone: M050-rfqm1p
key_files:
  - tests/integration/api_test.go
  - tests/integration/go.mod
  - benchmark-results/m050-s01-test-actuality.md
key_decisions:
  - Не использовать `FD_API_KEY` как fallback в root integration, чтобы не ловить stale shell secrets; требовать explicit `FD_INTEGRATION_API_KEY`.
duration: 
verification_result: passed
completed_at: 2026-06-15T14:40:09.677Z
blocker_discovered: false
---

# T03: Исправлен stale root integration test layer под текущий auth/runtime контракт.

**Исправлен stale root integration test layer под текущий auth/runtime контракт.**

## What Happened

Добавлен standalone module `tests/integration/go.mod`, обновлён `tests/integration/api_test.go`: base URL теперь задаётся через `FD_BASE_URL`, unauthenticated embeddings проверяет fail-closed 401, protected checks требуют явный `FD_INTEGRATION_API_KEY` и иначе skip. Это предотвращает ложные failures от случайного `FD_API_KEY` в shell и делает существующий integration layer runnable без чтения или печати секретов.

## Verification

Финальные проверки прошли: `cd api && go test ./...` — 295 passed; `cd tests/integration && go test -v .` — 2 passed in 1 package.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./...` | 0 | ✅ pass | 6900ms |
| 2 | `cd tests/integration && go test -v .` | 0 | ✅ pass | 6900ms |

## Deviations

Root integration happy-path с настоящим bearer token не выполнялся в S01, потому что matching running-service key не был доступен без чтения секрета; full authenticated live suite запланирован в S02.

## Known Issues

Protected root integration tests skip без `FD_INTEGRATION_API_KEY`.

## Files Created/Modified

- `tests/integration/api_test.go`
- `tests/integration/go.mod`
- `benchmark-results/m050-s01-test-actuality.md`
