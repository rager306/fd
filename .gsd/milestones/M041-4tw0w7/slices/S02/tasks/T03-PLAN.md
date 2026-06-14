---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Added /live and /ready probes backed by lifecycle state: /live always 200, /ready 200 after warmup and 503 model_not_loaded before warmup/shutdown.

api/handlers/probes.go: GET /live — cheap, проверяет только process alive, всегда 200 (даже если warmup not done). GET /ready — проверяет IsReady(), 200 если warmup done, 503 (overloaded_error, model_not_loaded, Retry-After: 5) если нет. Оба endpoints используют lifecycle state из T01.

## Inputs

- None specified.

## Expected Output

- `api/handlers/probes.go`
- `api/handlers/probes_test.go`

## Verification

Unit tests: после MarkWarmupDone → /ready 200; до → /ready 503 с code=model_not_loaded, Retry-After: 5. /live всегда 200.
