---
id: T02
parent: S01
milestone: M050-rfqm1p
key_files:
  - benchmark-results/m050-s01-test-actuality.md
key_decisions:
  - Использовать explicit `FD_INTEGRATION_API_KEY` для root integration, чтобы не зависеть от случайного shell `FD_API_KEY`.
duration: 
verification_result: mixed
completed_at: 2026-06-15T14:39:46.141Z
blocker_discovered: false
---

# T02: Запущены текущие test commands и выявлены устаревшие ожидания root integration layer.

**Запущены текущие test commands и выявлены устаревшие ожидания root integration layer.**

## What Happened

Регулярные `api` проверки (`go test`, `go test -short`, golangci-lint, govulncheck) прошли. Старый root command `go test ./tests/integration` не запускался из-за отсутствия root module. После первичной модульной попытки обнаружены две актуальности проблемы: `go 1.25` вместо `go 1.25.0` и зависимость protected checks от случайного `FD_API_KEY`, не совпадающего с running container.

## Verification

Команды и exit outcomes зафиксированы в `benchmark-results/m050-s01-test-actuality.md`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./...` | 0 | ✅ pass | 18200ms |
| 2 | `cd api && go test ./... -short` | 0 | ✅ pass | 18100ms |
| 3 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass | 18000ms |
| 4 | `cd api && go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 0 | ✅ pass | 18000ms |
| 5 | `go test ./tests/integration` | 1 | ❌ expected stale failure identified | 18000ms |

## Deviations

None.

## Known Issues

Protected happy-path root integration требует explicit `FD_INTEGRATION_API_KEY`; без него эти checks skip, а полный auth happy-path переносится в S02.

## Files Created/Modified

- `benchmark-results/m050-s01-test-actuality.md`
