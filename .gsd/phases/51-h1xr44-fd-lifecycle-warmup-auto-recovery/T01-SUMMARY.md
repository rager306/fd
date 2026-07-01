---
id: T01
parent: S01
milestone: M051-h1xr44
key_files:
  - (none)
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-07-01T09:17:34.679Z
blocker_discovered: false
---

# T01: Расширил стартовое warmup окно до 5×5с (~30s) через env-config FD_WARMUP_START_MAX_ATTEMPTS/FD_WARMUP_START_BACKOFF_SEC; main() переключён на env-based policy с логированием конфигурации.

**Расширил стартовое warmup окно до 5×5с (~30s) через env-config FD_WARMUP_START_MAX_ATTEMPTS/FD_WARMUP_START_BACKOFF_SEC; main() переключён на env-based policy с логированием конфигурации.**

## What Happened

Добавил `warmupRetryPolicyFromEnv()` в api/main.go рядом с `defaultWarmupRetryPolicy()`. Новая функция читает FD_WARMUP_START_MAX_ATTEMPTS (default 5) и FD_WARMUP_START_BACKOFF_SEC (default 5) через envutil.Int. Стартовое окно теперь ~30с (5 попыток × 5с backoff + per-call 5s timeout) против прежних ~6с (3 попытки × exponential 2s/4s), что покрывает TEI CPU BERT load ~15-20с. Backoff изменён с exponential на фиксированный: прежний doubling (2s,4s,8s) на больших попытках выглядел как hang и промахивался мимо типичного TEI load window. В main() вызов startModelWarmup(...) заменён на startModelWarmupWithPolicy(...) с env-policy и INFO-логом конфигурации. Сигнатуры warmupRetryPolicy, defaultWarmupRetryPolicy, startModelWarmup, startModelWarmupWithPolicy не изменены (LOW-impact по gitnexus). Существующие 4 теста TestStartModelWarmup* не сломаны: они передают policy явно и не зависят от env/default.

## Verification

go build ./api/... (exit 0); go test ./api/... -run TestStartModelWarmup -count=1 — все 4 существующих теста проходят (они передают warmupRetryPolicy явно в startModelWarmupWithPolicy, не зависят от env/default policy).

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go build ./...` | 0 | pass | 3200ms |
| 2 | `cd /root/fd/api && go test ./... -run TestStartModelWarmup -count=1 -timeout 30s` | 0 | pass | 750ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

None.
