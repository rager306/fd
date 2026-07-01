---
estimated_steps: 10
estimated_files: 1
skills_used: []
---

# T02: Implement startWarmupRecovery background loop in api/lifecycle/recovery.go

Создать новый файл `api/lifecycle/recovery.go` с функцией `startWarmupRecovery(ctx context.Context, logger *slog.Logger, state *lifecycle.State, model embed.Embedder, timeout, interval time.Duration, enabled bool)`. Поведение:
- Если enabled=false — no-op возврат (для feature-flag-off).
- Запустить goroutine, которая в select{ctx.Done | ticker}:
  - Если state.IsWarmupDone() или state.IsShuttingDown() — выход, log "recovery stopped".
  - Иначе PreWarm(ctx, model) с per-call timeout; при успехе — state.MarkWarmupDone(), exit; при ошибке — state.SetLastError(err), log WARN с attempt#, next_interval_sec.
- Ticker: `time.NewTicker(interval)`.
- Подсчёт попыток для логов через closure var.
- Возможный capped backoff (опционально): если число attempts > 10, увеличить interval пропорционально или перейти в quieter mode (например удвоенный interval). Простота важнее — первой версии достаточно фиксированного interval, без cap. Документировать это как известное ограничение, не блокер.

Не модифицировать `lifecycle.State.MarkWarmupDone/IsReady/SetLastError/IsShuttingDown` (HIGH-impact). Не модифицировать embed.Embedder interface.

Embedder цикл: model уже nil-guarded в PreWarm.

## Inputs

- `api/lifecycle/warmup.go`
- `api/lifecycle/state.go`
- `api/embed (для интерфейса)`

## Expected Output

- `api/lifecycle/recovery.go`

## Verification

go build ./api/lifecycle/... ; go vet ./api/lifecycle/...

## Observability Impact

Structured logs: warmup recovery attempt starting/warn, warmup recovery succeeded, warmup recovery stopped (reason). last_error обновляется в /health.
