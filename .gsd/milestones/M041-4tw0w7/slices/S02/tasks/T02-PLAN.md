---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Pre-warm model при старте

api/lifecycle/warmup.go: функция PreWarm(ctx, model) error которая вызывает model.Encode с 1 dummy input и логирует latency. Запускается из main.go после server start но ДО readiness=ready. Server start НЕ блокируется — http.Server.Serve() стартует сразу, lifecycle state warmupDone=false, /ready отвечает 503, /live отвечает 200 (process alive). Горутина warmup: loadModel → PreWarm → MarkWarmupDone. На любую ошибку — логируем, НО warmupDone остаётся false (server не ready, /ready=503, /health deep показывает error).

## Inputs

- None specified.

## Expected Output

- `api/lifecycle/warmup.go`
- `api/lifecycle/warmup_test.go`
- `api/main.go`

## Verification

Integration test: запустить fd binary, в течение 100ms после start curl /live → 200, curl /ready → 503, через 30s (после warmup) curl /ready → 200. Проверить логи: model loading, warmup started, warmup done.
