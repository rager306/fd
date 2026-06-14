# S02: Lifecycle warmup readiness and graceful shutdown

**Goal:** Lifecycle management: pre-warm model при старте, /live (cheap) и /ready (200 only after warmup), 503+Retry-After для model_not_loaded/model_overloaded/shutting_down, SIGTERM handler с in-flight drain за ≤30s. Закрывает R-P0-3, R-P0-4, R-P0-5.
**Demo:** After this, fd pre-warms model at startup, /live is cheap, /ready transitions 503 to 200 after warmup, shutdown gates new requests with 503+Retry-After, and the slice passes M043 gates: go test ./..., golangci-lint 18 linters, no reachable govulncheck findings.

## Must-Haves

- Startup sequence воспроизводится: /live 200 сразу, /ready 503 → 200 после warmup, /health deep status: down → ok
- F-1 scenario: caller hit во время warmup получает 503 model_not_loaded + Retry-After: 5
- F-2 scenario: caller hit при overload получает 503 model_overloaded + Retry-After: 5
- F-5 scenario: SIGTERM приводит к 503 shutting_down + Retry-After: 30 для новых запросов, in-flight завершается нормально, exit code 0
- Exit time после SIGTERM ≤ 35s (30s drain + 5s margin)
- Process exit 0 на чистом shutdown, exit 1 на force timeout
- golangci-lint pass

## Proof Level

- This slice proves: runtime + integration

## Integration Closure

Lifecycle state (warmupDone, shuttingDown, inflight) атомики используются в request handler (блокирует до warmup) и в S03 metrics gauge (model_loaded). Shutdown использует http.Server.Shutdown + WaitGroup.

## Verification

- Логируются lifecycle events: model loading, warmup started, warmup done, SIGTERM received, in-flight drained. S03 добавит fd_model_loaded gauge на этот state.

## Tasks

- [x] **T01: Added lifecycle State package with warmup/readiness/shutdown flags, in-flight request tracking, drain timeout, last error, and context helpers.** `est:2h`
  api/lifecycle/state.go: type State struct { warmupDone atomic.Bool; shuttingDown atomic.Bool; inflight sync.WaitGroup; lastError atomic.Value }. Methods: MarkWarmupDone(), IsReady() bool, BeginShutdown(), IsShuttingDown() bool, TrackRequest(start, done), WaitDrain(timeout) error. State синглтон, передаётся в handlers через context.
  - Files: `api/lifecycle/state.go`, `api/lifecycle/state_test.go`
  - Verify: Unit tests: MarkWarmupDone → IsReady true. BeginShutdown → IsShuttingDown true. TrackRequest properly tracks inflight, WaitDrain(0) returns immediately when empty, WaitDrain(100ms) blocks. Test concurrent shutdown while requests inflight.

- [x] **T02: Added async model pre-warm: server startup no longer blocks on warmup, lifecycle state flips ready only after successful dummy embedding.** `est:3h`
  api/lifecycle/warmup.go: функция PreWarm(ctx, model) error которая вызывает model.Encode с 1 dummy input и логирует latency. Запускается из main.go после server start но ДО readiness=ready. Server start НЕ блокируется — http.Server.Serve() стартует сразу, lifecycle state warmupDone=false, /ready отвечает 503, /live отвечает 200 (process alive). Горутина warmup: loadModel → PreWarm → MarkWarmupDone. На любую ошибку — логируем, НО warmupDone остаётся false (server не ready, /ready=503, /health deep показывает error).
  - Files: `api/lifecycle/warmup.go`, `api/lifecycle/warmup_test.go`, `api/main.go`
  - Verify: Integration test: запустить fd binary, в течение 100ms после start curl /live → 200, curl /ready → 503, через 30s (после warmup) curl /ready → 200. Проверить логи: model loading, warmup started, warmup done.

- [x] **T03: Added /live and /ready probes backed by lifecycle state: /live always 200, /ready 200 after warmup and 503 model_not_loaded before warmup/shutdown.** `est:1h`
  api/handlers/probes.go: GET /live — cheap, проверяет только process alive, всегда 200 (даже если warmup not done). GET /ready — проверяет IsReady(), 200 если warmup done, 503 (overloaded_error, model_not_loaded, Retry-After: 5) если нет. Оба endpoints используют lifecycle state из T01.
  - Files: `api/handlers/probes.go`, `api/handlers/probes_test.go`
  - Verify: Unit tests: после MarkWarmupDone → /ready 200; до → /ready 503 с code=model_not_loaded, Retry-After: 5. /live всегда 200.

- [x] **T04: Added lifecycle gate middleware for /v1/embeddings with warmup/shutdown rejection and in-flight request tracking.** `est:2h`
  api/middleware/lifecycle.go: gin middleware который проверяет IsReady() и !IsShuttingDown() перед передачей в handler. Если !IsReady → 503 model_not_loaded + Retry-After: 5. Если IsShuttingDown → 503 shutting_down + Retry-After: 30. Также TrackRequest(start, done) для inflight tracking. Подключается в router setup после validation (S01), до embed handler.
  - Files: `api/middleware/lifecycle.go`, `api/middleware/lifecycle_test.go`
  - Verify: Unit tests: до warmup → 503 model_not_loaded, Retry-After: 5. После BeginShutdown → 503 shutting_down, Retry-After: 30. inflight counter инкрементируется на request и декрементируется на response.

- [ ] **T05: Graceful shutdown по SIGTERM/SIGINT** `est:3h`
  api/lifecycle/shutdown.go: signal handler для SIGTERM и SIGINT. По сигналу: BeginShutdown(), log SIGTERM received, http.Server.Shutdown(ctxWith30sTimeout) — отказывает в новых соединениях, ждёт активные handlers до 30s, force close после. WaitDrain(30s) с inflight tracking. Exit 0 на clean drain, exit 1 на force timeout. Также: при shutdown in-flight handlers получают 503 shutting_down+Retry-After: 30 (через lifecycle middleware из T04).
  - Files: `api/lifecycle/shutdown.go`, `api/lifecycle/shutdown_test.go`, `api/main.go`
  - Verify: Integration test: запустить fd, послать 1 long-running request, послать SIGTERM, новые запросы получают 503 shutting_down+Retry-After: 30, in-flight завершается нормально, process exit 0, total time ≤ 35s. Также test: SIGTERM при idle → exit < 1s.

- [ ] **T06: Integration tests для behavior scenarios F-1/F-2/F-5** `est:3h`
  tests/integration/fd_v2_lifecycle_test.go: воспроизвести F-1 (caller hit во время warmup → 503 model_not_loaded + Retry-After), F-2 (concurrent overload → 503 model_overloaded + Retry-After, после снижения load → 200), F-5 (SIGTERM → 503 shutting_down + drain). Также test: startup sequence — /live=200, /ready=503, /ready=200 после warmup, /health deep корректно меняется. Спека: docs/fd-v2.md Section 6.1 + 6.3 F-1/F-2/F-5.
  - Files: `tests/integration/fd_v2_lifecycle_test.go`
  - Verify: go test ./tests/integration/... -run TestFdV2Lifecycle -v: F-1, F-2, F-5, и startup sequence test все pass.

## Files Likely Touched

- api/lifecycle/state.go
- api/lifecycle/state_test.go
- api/lifecycle/warmup.go
- api/lifecycle/warmup_test.go
- api/main.go
- api/handlers/probes.go
- api/handlers/probes_test.go
- api/middleware/lifecycle.go
- api/middleware/lifecycle_test.go
- api/lifecycle/shutdown.go
- api/lifecycle/shutdown_test.go
- tests/integration/fd_v2_lifecycle_test.go
