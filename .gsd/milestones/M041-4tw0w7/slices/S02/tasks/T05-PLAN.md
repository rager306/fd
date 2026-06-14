---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T05: Graceful shutdown по SIGTERM/SIGINT

api/lifecycle/shutdown.go: signal handler для SIGTERM и SIGINT. По сигналу: BeginShutdown(), log SIGTERM received, http.Server.Shutdown(ctxWith30sTimeout) — отказывает в новых соединениях, ждёт активные handlers до 30s, force close после. WaitDrain(30s) с inflight tracking. Exit 0 на clean drain, exit 1 на force timeout. Также: при shutdown in-flight handlers получают 503 shutting_down+Retry-After: 30 (через lifecycle middleware из T04).

## Inputs

- None specified.

## Expected Output

- `api/lifecycle/shutdown.go`
- `api/lifecycle/shutdown_test.go`
- `api/main.go`

## Verification

Integration test: запустить fd, послать 1 long-running request, послать SIGTERM, новые запросы получают 503 shutting_down+Retry-After: 30, in-flight завершается нормально, process exit 0, total time ≤ 35s. Также test: SIGTERM при idle → exit < 1s.
