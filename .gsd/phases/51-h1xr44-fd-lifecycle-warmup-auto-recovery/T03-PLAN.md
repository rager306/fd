---
estimated_steps: 6
estimated_files: 1
skills_used: []
---

# T03: Wire startWarmupRecovery into main() with env-config and shutdown context

Спланировать запуск recovery из main(): вызов `startWarmupRecovery(ctx, logger, lifecycleState, embeddingClient, defaultWarmupTimeout, recoveryInterval, recoveryEnabled)` ПОСЛЕ `startModelWarmup(...)`.

Контекст: создать `ctx, cancel := context.WithCancel(context.Background())`, передать в recovery; при получении сигнала (`AwaitSignalAndShutdown`) вызвать cancel() для немедленного выхода recovery goroutine. Способ привязки:
- Либо обернуть signal-handling: сделать так, чтобы перед AwaitSignalAndShutdown cancel вызывался по signalCh (spawn small goroutine или использовать ctx.Done в recovery через переданный ctx).
- Простейший вариант: recovery использует `signalContext` (ctx, cancel) который cancel'ится первым в сигнальном flow. Если lifecycle BeginShutdown уже вызывается в GracefulShutdown, recovery сам проверяет state.IsShuttingDown() и выходит — cancel не строго нужен, но добавляет мгновенный exit без ожидания следующего tick.

env-config через envutil.Int: `FD_WARMUP_RECOVERY_INTERVAL_SEC` (default 30), `FD_WARMUP_RECOVERY_ENABLED` (env bool: default true через truthy check). Использовать существующий паттерн env parsing в main.go.

Логировать старта: `logger.Info("warmup recovery enabled", "interval", interval, "timeout", timeout)`.

## Inputs

- `api/main.go`
- `api/lifecycle/recovery.go`
- `api/internal/envutil/int.go`

## Expected Output

- `api/main.go (только wiring в main функции)`

## Verification

go build ./api/... ; go vet ./api/...

## Observability Impact

INFO log "warmup recovery enabled" с interval/timeout. INFO log "warmup recovery stopped reason=shutdown/success".
