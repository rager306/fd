---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Endpoints /version, /info, /v1/healthcheck

api/handlers/observability.go: GET /version — возвращает buildinfo.Info + uptime. GET /info — возвращает список моделей с dims=[512,1024], max_input_length_tokens=512, max_batch_size=32, loaded, warmup_done, device (cuda:0/cpu). GET /v1/healthcheck — alias для /health, тот же response. Все endpoints используют lifecycle state из S02.

## Inputs

- None specified.

## Expected Output

- `api/handlers/observability.go`
- `api/handlers/observability_test.go`

## Verification

Integration tests: T-H-10 (/version 200 с version field), T-H-7 (/health deep с model_loaded, warmup_done), T-E-1..T-E-3 (Section 5.5 existence: /version 200, /info 200, /metrics 200 text/plain).
