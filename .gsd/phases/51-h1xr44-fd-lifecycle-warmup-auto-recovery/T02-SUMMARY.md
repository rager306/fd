---
id: T02
parent: S01
milestone: M051-h1xr44
key_files:
  - api/lifecycle/recovery.go
key_decisions:
  - Не модифицирую lifecycle.State — recovery использует только публичный API (IsWarmupDone, IsShuttingDown, SetLastError, MarkWarmupDone)
  - Goroutine exits на терминальных условиях — fixed interval, без cap в первой версии
  - Per-call timeout через context.WithTimeout — медленный TEI не пиннит goroutine
duration: 
verification_result: passed
completed_at: 2026-07-01T09:18:24.257Z
blocker_discovered: false
---

# T02: Реализована startWarmupRecovery в новом api/lifecycle/recovery.go: background goroutine с periodic retry, exits на IsWarmupDone/IsShuttingDown/ctx.Done, per-call timeout, structured WARN/INFO logs

**Реализована startWarmupRecovery в новом api/lifecycle/recovery.go: background goroutine с periodic retry, exits на IsWarmupDone/IsShuttingDown/ctx.Done, per-call timeout, structured WARN/INFO logs**

## What Happened

Создан новый файл api/lifecycle/recovery.go с экспортируемой функцией StartWarmupRecovery. Реализация: 
- На enabled=false/interval<=0/nil-state/model — no-op или WARN и return.
- goroutine в select{ctx.Done | ticker.C}: проверяет state.IsWarmupDone и state.IsShuttingDown перед каждой попыткой (защита от race между tick и проверкой), делает PreWarm с per-call timeout через context.WithTimeout.
- Успех → state.MarkWarmupDone() + INFO log.
- Ошибка → state.SetLastError(err) + WARN log с attempt# и next_interval.
- Exits deterministically на: ctx.Done, state.IsWarmupDone (уже или во время тика), state.IsShuttingDown, успех.

Дизайн-решения:
1. Не модифицирую lifecycle.State — recovery использует только публичный API. Это устраняет HIGH-risk impact на MarkWarmupDone (gitnexus показал 2 прямых caller'а: startModelWarmupWithPolicy и WarmupHandler.runWarmup).
2. Не координирую с WarmupHandler.inProgress — PreWarm idempotent, double-count в логах приемлем как known issue. Будет работать, но если будет нужен — отдельная T0X.
3. Без capped backoff в первой версии — комментарий в коде документирует это как осознанное упрощение.

Не использованы никакие новые методы State — interface контракт не изменён. embed.Embedder interface не тронут.

## Verification

go build ./api/lifecycle/... (exit 0); go vet ./api/lifecycle/... (exit 0). Реализация StartWarmupRecovery в api/lifecycle/recovery.go использует только уже существующие публичные методы State — без модификации State (HIGH-impact по gitnexus для MarkWarmupDone). embed.Embedder interface не тронут.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go build ./lifecycle/...` | 0 | pass | 2100ms |
| 2 | `cd /root/fd/api && go vet ./lifecycle/...` | 0 | pass | 1600ms |

## Deviations

Добавил защиту interval <= 0: WARN лог + return вместо panic — main() может передать некорректный interval при env misconfiguration.

## Known Issues

Log-spam при permanent TEI unavailability: документировано, deferred до T0X-ного регламента по нагрузке в проде.

## Files Created/Modified

- `api/lifecycle/recovery.go`
