---
id: S01
parent: M051-h1xr44
milestone: M051-h1xr44
provides:
  - (none)
requires:
  []
affects:
  []
key_files:
  - api/main.go
  - api/lifecycle/recovery.go
  - api/lifecycle/recovery_test.go
  - api/internal/envutil/bool.go
  - api/internal/envutil/bool_test.go
  - tools/verify_warmup_recovery.sh
  - .gsd/runtime/M051-h1xr44/api-recovery-logs.jsonl
  - .gsd/runtime/M051-h1xr44/warmup-recovery-runtime-trace.log
key_decisions:
  - D052: combined approach — extended startup window (5×5s) + periodic recovery (30s interval)
  - Backoff изменён с exponential на fixed — doubling промахивается мимо TEI load window
  - Recovery НЕ координируется с WarmupHandler.inProgress — PreWarm idempotent
  - Без capped backoff в v1 — deferred до real log-spam evidence
  - All env knobs backward-compatible с safe defaults, feature-flag rollback via FD_WARMUP_RECOVERY_ENABLED=false
patterns_established:
  - envutil.BoolOrDefault для bool env parsing (truthy/falsy/unknown→fallback, безопасный misconfig)
  - lifecycle recovery pattern: periodic background goroutine with select{ctx.Done | ticker} exits on terminal conditions
  - Separation: startup warmup (bounded attempts) vs recovery (periodic until terminal) — different policies, same State API
observability_surfaces:
  - none
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-07-01T10:50:21.941Z
blocker_discovered: false
---

# S01: Warmup recovery contract

**fd-api восстанавливает readiness автоматически: расширенное стартовое окно 5×5с + periodic recovery каждые 30с с feature flag; end-to-end подтверждено — recovery succeeded attempt:7 после 3.5 мин TEI CPU load.**

## What Happened

M051-h1xr44 S01 реализован в 5 задачах:

T01 — api/main.go: warmupRetryPolicyFromEnv() читает FD_WARMUP_START_MAX_ATTEMPTS (default 5, было 3) и FD_WARMUP_START_BACKOFF_SEC (default 5, было exponential 2s<<n). Backoff изменён с exponential на фиксированный — doubling 2s/4s/8s промахивался мимо TEI load window и выглядел как hang. Стартовое окно: ~30с вместо ~6с.

T02 — api/lifecycle/recovery.go (новый файл): StartWarmupRecovery goroutine с periodic retry. Exits на IsWarmupDone/IsShuttingDown/ctx.Done. Per-call timeout через context.WithTimeout — медленный TEI не пиннит goroutine. Structured INFO/WARN логи. НЕ модифицирует lifecycle.State (HIGH-impact по gitnexus на MarkWarmupDone — 2 direct callers: startModelWarmupWithPolicy и WarmupHandler.runWarmup).

T03 — wiring в main(): recoveryCtx + defer cancel, env-knobs FD_WARMUP_RECOVERY_INTERVAL_SEC (default 30s), FD_WARMUP_RECOVERY_ENABLED (default true, opt-out). Новый envutil.BoolOrDefault helper (api/internal/envutil/bool.go) с 5 тестами — truthy/falsy/unknown→fallback, whitespace-safe.

T04 — api/lifecycle/recovery_test.go: 8 тестов покрывают state machine (disabled/already-warm/retry-until-success/stop-at-success/stop-at-shutdown/stop-at-ctx-cancel/nil-args/zero-interval). Все проходят с -race. Весь api test suite (11 пакетов) зелёный с -race.

T05 — tools/verify_warmup_recovery.sh: docker compose reproducer race condition. End-to-end recovery подтверждён дважды в живом окружении: startup 5 attempts failed (20с) → recovery attempt 1-6 every 30с → TEI Ready (3 мин CPU load) → recovery succeeded attempt:7 latency:197ms elapsed:210197ms → /health ok. Evidence в .gsd/runtime/M051-h1xr44/.

Архитектурные решения:
- D052 зафиксировал выбор combined подхода (extended startup + periodic recovery)
- Recovery goroutine не координируется с WarmupHandler.inProgress — PreWarm idempotent, double-count в логах acceptable
- Без capped backoff в первой версии — документировано, deferred до real log-spam evidence
- env knobs все с safe defaults, backward-compatible, feature-flag rollback через FD_WARMUP_RECOVERY_ENABLED=false

Out-of-scope (пойдёт в отдельный milestone — волна 3): env-file аномалия docker compose up --force-recreate не инжектит FD_API_KEY; не блокирует recovery contract, блокирует только smoke /v1/embeddings.

## Verification

go test ./api/... -race -count=1 — all 11 packages green. Recovery confirmed end-to-end twice via docker logs + /health: startup warmup 5 attempts fail → recovery 6 attempts fail (30s interval exact) → TEI ready → recovery attempt 7 succeeded → /health ok. Evidence: .gsd/runtime/M051-h1xr44/api-recovery-logs.jsonl (17 lines) + warmup-recovery-runtime-trace.log (87 lines).

## Requirements Advanced

- R-new — T01-T05 реализуют и доказывают auto-recovery contract

## Requirements Validated

None.

## New Requirements Surfaced

- R-new: fd-api must automatically recover model readiness after transient startup-time TEI unreachability without operator intervention

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

Recovery без capped backoff — при permanent TEI outage WARN логи шумят каждые 30с. Документировано как осознанное упрощение; cap добавляется отдельной T0X если log-spam станет real problem в проде.

## Follow-ups

Отдельный milestone (волна 3): env-file аномалия docker compose up --force-recreate не инжектит FD_API_KEY в контейнер — блокирует smoke /v1/embeddings, но не recovery contract.

## Files Created/Modified

- `api/main.go` — warmupRetryPolicyFromEnv() + wiring startWarmupRecovery в main()
- `api/lifecycle/recovery.go` — новый файл: StartWarmupRecovery background goroutine
- `api/lifecycle/recovery_test.go` — новый файл: 8 recovery state machine тестов
- `api/internal/envutil/bool.go` — новый файл: BoolOrDefault helper
- `api/internal/envutil/bool_test.go` — новый файл: 5 BoolOrDefault тестов
- `tools/verify_warmup_recovery.sh` — новый файл: docker compose reproducer
