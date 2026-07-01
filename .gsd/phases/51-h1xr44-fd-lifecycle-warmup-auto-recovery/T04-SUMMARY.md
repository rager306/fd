---
id: T04
parent: S01
milestone: M051-h1xr44
key_files:
  - (none)
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-07-01T09:25:41.808Z
blocker_discovered: false
---

# T04: 8 recovery-тестов в api/lifecycle/recovery_test.go покрывают disabled/already-warm/retry-until-success/stop-at-success/stop-at-shutdown/stop-at-ctx-cancel/nil-args/zero-interval; все проходят с -race, весь api-сьют зелёный

**8 recovery-тестов в api/lifecycle/recovery_test.go покрывают disabled/already-warm/retry-until-success/stop-at-success/stop-at-shutdown/stop-at-ctx-cancel/nil-args/zero-interval; все проходят с -race, весь api-сьют зелёный**

## What Happened

Создан api/lifecycle/recovery_test.go с 8 тестами, покрывающими recovery state machine. Тесты переиспользуют существующие в пакете хелперы warmupModelFunc (warmup_test.go) и discardLogger (shutdown_test.go); waitForCondition добавлен локально (2с deadline вместо 1с в main_test.go для стабильности на нагруженном CI).

Покрытие:
1. DisabledIsNoop — enabled=false: модель не вызывается.
2. AlreadyWarmNoop — state.MarkWarmupDone() до старта: модель не вызывается, immediate exit.
3. RetriesAndMarksReady — boom×2, успех×3: state.IsReady() в конце, attempts=3, LastError=nil.
4. StopsAtSuccess — успех с первой попытки: attempts=1, goroutine не продолжает после первого tick.
5. StopsAtShutdown — модель всегда boom, BeginShutdown через 40ms: рост attempts останавливается.
6. StopsAtContextCancel — модель всегда boom, ctx.cancel через 40ms: рост останавливается.
7. NilArgsDontPanic — nil state/model не паникуют, no-op.
8. NonPositiveIntervalNoop — interval=0: WARN и return без goroutine.

Все тесты проходят с -race. Весь api-тест-сьют (11 пакетов) зелёный с -race.

Исправил один баг в тесте: в StopsAtSuccess модель изначально не инкрементировала attempts counter — это был мой промах в тесте (не баг в коде). После добавления attempts.Add(1) тест проходит.

## Verification

go test ./api/lifecycle/... -run TestStartWarmupRecovery -count=1 -race (pass, 8 тестов); go test ./api/... -count=1 -race (pass, все 11 пакетов). Race detector чист.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go test ./lifecycle/... -run TestStartWarmupRecovery -count=1 -race -timeout 60s` | 0 | pass | 1553ms |
| 2 | `cd /root/fd/api && go test ./... -count=1 -race -timeout 120s` | 0 | pass | 13000ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

None.
