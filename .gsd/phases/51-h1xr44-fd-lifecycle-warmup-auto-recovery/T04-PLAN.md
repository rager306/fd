---
estimated_steps: 8
estimated_files: 1
skills_used: []
---

# T04: Unit tests for recovery state machine

Добавить unit-тесты для recovery state machine в `api/lifecycle/recovery_test.go` (по образцу `api/main_test.go:warmupModelFunc`). Тесты:

1. `TestStartWarmupRecoveryRetriesAndMarksReady` — модель возвращает boom первые 2 раза, успех на 3-й. interval=10ms, timeout=1s. Assert: state.IsReady() true в конце, SetLastError=nil, ровно 3 PreWarm вызова.
2. `TestStartWarmupRecoveryStopsAtShutdown` — модель всегда возвращает boom, но через 30ms вызываем state.BeginShutdown(). Assert: recovery goroutine завершилась, нет утечки, число attempts ограничено (например ≤3 за 50ms окно).
3. `TestStartWarmupRecoveryStopsAtSuccess` — модель возвращает успех с первой попытки. Recovery должен вызвать MarkWarmupDone и выйти. Assert: только 1 вызов PreWarm, state.IsReady.
4. `TestStartWarmupRecoveryDisabledIsNoop` — enabled=false. Assert: goroutine не запущена, модель никогда не вызвана.
5. `TestStartWarmupRecoveryAlreadyWarmNoop` — state.MarkWarmupDone() вызван до запуска recovery. Assert: модель не вызвана, immediate exit.

Использовать pattern `waitForCondition` из main_test.go (или перенести хелпер). Интервалы тестовые: 10-20ms чтобы тесты были быстрые (<1s).

Документировать в комментарии, что interval=10ms НЕ прод-настройка — это тестовый шим для скорости.

## Inputs

- `api/lifecycle/recovery.go`
- `api/main_test.go (warmupModelFunc pattern)`

## Expected Output

- `api/lifecycle/recovery_test.go`

## Verification

go test ./api/lifecycle/... -run TestStartWarmupRecovery -count=1 -race

## Observability Impact

Нет — unit tests.
