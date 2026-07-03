---
id: T01
parent: S03
milestone: M050-rfqm1p
key_files:
  - benchmark-results/m050-s03-mutation-baseline.md
key_decisions:
  - Use avito-tech go-mutesting fork for M050 bounded baseline.
duration: 
verification_result: passed
completed_at: 2026-06-15T14:54:16.557Z
blocker_discovered: false
---

# T01: Выбран рабочий Go mutation runner и подтверждён smoke-запуск.

**Выбран рабочий Go mutation runner и подтверждён smoke-запуск.**

## What Happened

Проверены `zimmski/go-mutesting` и более свежий `avito-tech/go-mutesting`. Для baseline выбран `github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest`, потому что он работает с Go 1.25 toolchain. Smoke на `api/cache/hash.go` прошёл со score 1.0: 2 mutants killed, 0 survived.

## Verification

Smoke command `cd api && go run github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest --exec 'go test ./cache' --exec-timeout 30 ./cache/hash.go` exited 0 with mutation score 1.000000.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go run github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest --exec 'go test ./cache' --exec-timeout 30 ./cache/hash.go` | 0 | ✅ pass: mutation score 1.000000 (2 killed, 0 survived) | 6200ms |

## Deviations

None.

## Known Issues

Runner requires Go >= 1.25.5 and auto-switches to go1.25.11 in this environment.

## Files Created/Modified

- `benchmark-results/m050-s03-mutation-baseline.md`
