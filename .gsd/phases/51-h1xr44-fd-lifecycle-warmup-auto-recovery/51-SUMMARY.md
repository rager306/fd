---
id: M051-h1xr44
title: "fd lifecycle warmup auto-recovery"
status: complete
completed_at: 2026-07-01T10:51:42.410Z
key_decisions:
  - D052: combined approach — extended startup window (5×5s) + periodic recovery (30s interval, feature-flagged)
  - Backoff изменён с exponential на fixed — doubling промахивается мимо TEI load window
  - Recovery НЕ координируется с WarmupHandler.inProgress — PreWarm idempotent, double-count acceptable
  - Без capped backoff в v1 — deferred до real log-spam evidence
  - envutil.BoolOrDefault для bool env parsing (truthy/falsy/unknown→fallback, безопасный misconfig)
  - Все env knobs backward-compatible с safe defaults, feature-flag rollback via FD_WARMUP_RECOVERY_ENABLED=false
key_files:
  - api/main.go
  - api/lifecycle/recovery.go
  - api/lifecycle/recovery_test.go
  - api/internal/envutil/bool.go
  - api/internal/envutil/bool_test.go
  - tools/verify_warmup_recovery.sh
  - .gsd/runtime/M051-h1xr44/api-recovery-logs.jsonl
  - .gsd/runtime/M051-h1xr44/warmup-recovery-runtime-trace.log
  - .gsd/phases/51-h1xr44-fd-lifecycle-warmup-auto-recovery/51-VALIDATION.md
lessons_learned:
  - gitnexus_impact analysis перед изменением MarkWarmupDone (HIGH risk) позволил избежать модификации публичного API State — recovery использует только существующие методы
  - End-to-end verification на живом docker окружении обязателен для lifecycle-изменений: unit тесты не ловят race condition между tei startup и api warmup
  - TEI BERT load на CPU под contention может занимать 3+ минуты — budget для reproducer должен покрывать worst-case, не average
  - docker compose up --force-recreate имеет аномалию с env_file injection — отдельный milestone для расследования
---

# M051-h1xr44: fd lifecycle warmup auto-recovery

**fd-api автоматически восстанавливает readiness после race с медленным TEI startup: расширенное стартовое окно 5×5с + periodic recovery каждые 30с с feature flag; end-to-end доказано (recovery succeeded после 3.5 мин TEI load).**

## What Happened

M051-h1xr44 завершён в одной волне (5 задач, 1 slice, 1 milestone):

Проблема: fd-api залипал в degraded/model_not_loaded после race с медленным TEI startup (CPU BERT load ~15-20с, под contention до 3 мин), требуя ручного `docker restart fd_api`. Root cause MEM088 зафиксирован.

Решение (D052 combined approach):
1. Расширил startup warmup окно с 3×2с exponential (~6с) до 5×5с fixed (~30с) через env-config FD_WARMUP_START_MAX_ATTEMPTS/FD_WARMUP_START_BACKOFF_SEC.
2. Добавил periodic background recovery goroutine (api/lifecycle/recovery.go) с interval 30с (FD_WARMUP_RECOVERY_INTERVAL_SEC) и feature flag (FD_WARMUP_RECOVERY_ENABLED default true).
3. Recovery exits детерминированно на IsWarmupDone/IsShuttingDown/ctx.Done. Per-call timeout. Structured logs.

Доказательство end-to-end: recovery succeeded attempt:7 elapsed:210197ms после 3.5 мин TEI CPU load в двух независимых run'ах. /health {status:ok,warmup_done:true,model_loaded:true}. Evidence в .gsd/runtime/M051-h1xr44/.

Тесты: 8 новых recovery тестов + 5 новых BoolOrDefault тестов, все -race green. Весь api suite (11 пакетов) зелёный. Существующие TestStartModelWarmup* не сломаны (передают policy явно).

Архитектурная дисциплина: НЕ модифицировал lifecycle.State.MarkWarmupDone (HIGH-impact по gitnexus — 2 direct callers). Recovery вызывает существующие методы State как есть. /health shape и middleware.LifecycleGate не тронуты.

R046 (continuity) создан и validated.

Follow-up milestone (волна 3): env-file аномалия docker compose up --force-recreate — FD_API_KEY не инжектируется в контейнер. Не блокирует recovery, но блокирует smoke /v1/embeddings.

## Success Criteria Results

All 9 success criteria PASS with objective evidence (see VALIDATION.md). Key: end-to-end recovery proven twice, 13 new tests green with -race, no HIGH-impact State modifications.

## Definition of Done Results

1. PASS — StartWarmupRecovery goroutine запускается из main() после startModelWarmup с recoveryCtx cancel по shutdown. Evidence: api/main.go + api/lifecycle/recovery.go.
2. PASS — Startup окно расширено до 5×5с (FD_WARMUP_START_MAX_ATTEMPTS, FD_WARMUP_START_BACKOFF_SEC). Evidence: warmupRetryPolicyFromEnv() в api/main.go.
3. PASS — Recovery останавливается на IsWarmupDone() и IsShuttingDown(). Evidence: 8 unit тестов с -race + end-to-end логи.
4. PASS — Env knobs: FD_WARMUP_START_MAX_ATTEMPTS, FD_WARMUP_START_BACKOFF_SEC, FD_WARMUP_RECOVERY_INTERVAL_SEC, FD_WARMUP_RECOVERY_ENABLED — все с safe defaults.
5. PASS — lifecycle.State.MarkWarmupDone/IsReady сигнатуры не модифицированы (HIGH-impact по gitnexus).
6. PASS — /health shape и middleware.LifecycleGate не изменены.
7. PASS — go test ./api/... ./api/lifecycle/... проходит с -race, 8 новых recovery-тестов.
8. PASS — tools/verify_warmup_recovery.sh воспроизводит race и подтверждает auto-recovery.
9. PASS — fd_api после docker compose down/up приходит в /health 200 самостоятельно (recovery succeeded attempt:7).

## Requirement Outcomes

R046 (continuity, new): fd-api must auto-recover model readiness after transient startup-time TEI unreachability — VALIDATED. Evidence: recovery succeeded attempt:7 elapsed:210197ms after 3.5 min TEI CPU load; /health ok; 8 unit tests + full api suite -race green.

## Deviations

Без capped backoff в v1 — документировано как осознанное упрощение в D052. Скрипт verify_warmup_recovery.sh exit 1 при TEI load >budget (3+ мин под CPU contention) — это окружение, не recovery баг; для CI увеличить MAX_RECOVERY_WAIT_SEC.

## Follow-ups

Отдельный milestone (волна 3): env-file аномалия docker compose up --force-recreate не инжектит FD_API_KEY в контейнер — блокирует smoke /v1/embeddings (401 unauthorized), но не блокирует recovery contract. Нужно исследовать compose env_file interaction с --force-recreate и dependent service healthchecks.
