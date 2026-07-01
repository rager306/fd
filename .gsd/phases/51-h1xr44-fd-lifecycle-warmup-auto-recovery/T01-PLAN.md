---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Extend startup warmup retry policy with env-configurable attempts/backoff

Расширить existing `defaultWarmupRetryPolicy()` в `api/main.go` так, чтобы maxAttempts и backoff брались из env через `api/internal/envutil.Int`. Дефолты: maxAttempts=5 (текущее значение 3), backoff=5s (текущее значение 2s << (attempt-1) — заменить на фиксированный). Env knobs: `FD_WARMUP_START_MAX_ATTEMPTS` (default 5), `FD_WARMUP_START_BACKOFF_SEC` (default 5). Логика: стартовое окно должно быть достаточно чтобы покрыть ~15-20с TEI CPU load. Поведение существующих тестов `TestStartModelWarmupRetriesAndMarksReady` и `TestStartModelWarmupRecordsTerminalErrorAfterMaxAttempts` не должно сломаться — они передают `warmupRetryPolicy` явно в `startModelWarmupWithPolicy`, поэтому дефолтная политика может меняться свободно. Не модифицировать сигнатуру `warmupRetryPolicy` и сигнатуру `startModelWarmupWithPolicy` — они тестируются. Не трогать `lifecycle.State.MarkWarmupDone` (HIGH-impact по gitnexus).

## Inputs

- `api/main.go`
- `api/main_test.go`
- `api/internal/envutil/int.go`

## Expected Output

- `api/main.go (только defaultWarmupRetryPolicy и её вызовы)`

## Verification

go test ./api/... -run TestStartModelWarmup -count=1 (4 существующих теста должны проходить)

## Observability Impact

INFO/WARN логи остаются как есть; только max_attempts/logged будет новое значение.
