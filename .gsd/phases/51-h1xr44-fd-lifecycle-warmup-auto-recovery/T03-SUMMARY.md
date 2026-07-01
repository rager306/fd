---
id: T03
parent: S01
milestone: M051-h1xr44
key_files:
  - api/main.go
  - api/internal/envutil/bool.go
  - api/internal/envutil/bool_test.go
key_decisions:
  - recoveryCancel в defer — гарантирует остановку recovery даже если GracefulShutdown паникует
  - Env-knob FD_WARMUP_RECOVERY_ENABLED = true по умолчанию (opt-out для отката)
  - Env-knob FD_WARMUP_RECOVERY_INTERVAL_SEC = 30s по умолчанию (тише, чем у startup)
  - BoolOrDefault через envutil (а не локальный helper) — переиспользуемо для будущих bool-фич-флагов fd
  - BoolOrDefault unknown value → fallback (не паника, не silent true) — safer for misconfig
duration: 
verification_result: passed
completed_at: 2026-07-01T09:19:57.753Z
blocker_discovered: false
---

# T03: Recovery подключён в main(): recoveryCtx + cancel в defer, env-knobs FD_WARMUP_RECOVERY_INTERVAL_SEC (default 30s) и FD_WARMUP_RECOVERY_ENABLED (default true); envutil.BoolOrDefault новый helper с 5 тестами парсинга.

**Recovery подключён в main(): recoveryCtx + cancel в defer, env-knobs FD_WARMUP_RECOVERY_INTERVAL_SEC (default 30s) и FD_WARMUP_RECOVERY_ENABLED (default true); envutil.BoolOrDefault новый helper с 5 тестами парсинга.**

## What Happened

Добавлен wiring recovery после startModelWarmup. Создан recoveryCtx через context.WithCancel, defer cancel() обеспечивает cleanup даже при панике. Env-knobs:
- FD_WARMUP_RECOVERY_INTERVAL_SEC (default 30) → time.Duration через envutil.Int.
- FD_WARMUP_RECOVERY_ENABLED (default true) → bool через новый envutil.BoolOrDefault.

Логируем конфигурацию INFO'ом "warmup recovery config" с обоими полями перед StartWarmupRecovery.

В api/internal/envutil/bool.go добавлен BoolOrDefault: truthy = {"1","true","yes","on","y","t"} (case-insensitive), falsy = {"0","false","no","off","n","f"}, пустое или неизвестное → fallback (не silent toggle, безопасный misconfig fallback). Тесты покрывают все ветви (5 шт).

Сигнальный flow не изменён: recovery goroutine проверяет state.IsShuttingDown(), который выставляется в lifecycle.GracefulShutdown перед стартом drain. defer recoveryCancel() страхует медленный случай — recovery выйдет на следующем tick (≤ interval) либо мгновенно на ctx.Done.

Public API lifecycle.StartWarmupRecovery уже использовался как есть — изменения только в api/main.go (wiring) и envutil (новый helper). State, middleware, /health — не тронуты.

## Verification

go build ./api/... (exit 0); go test ./api/internal/envutil/... (pass, 5 BoolOrDefault тестов); существующие TestStartModelWarmup* не задеты.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go build ./...` | 0 | pass | 2900ms |
| 2 | `cd /root/fd/api && go test ./internal/envutil/... -count=1 -timeout 30s` | 0 | pass | 50ms |

## Deviations

Возможная конкуренция cancel + shutdown signal handling: при SIGTERM GracefulShutdown вызывает state.BeginShutdown() первым; recovery goroutine проверяет IsShuttingDown и выходит (см. recovery.go). cancel() в defer страхует медленный случай.

## Known Issues

recoveryCancel в defer сработает после main() return, но goroutine recovery уже могла завершиться через IsShuttingDown(), который выставляется в GracefulShutdown — double coverage OK.

## Files Created/Modified

- `api/main.go`
- `api/internal/envutil/bool.go`
- `api/internal/envutil/bool_test.go`
