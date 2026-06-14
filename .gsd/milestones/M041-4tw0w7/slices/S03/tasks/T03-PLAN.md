---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Deep /health с status/degraded/down

api/handlers/health.go (replace existing): GET /health возвращает { status: ok|degraded|down, time, model_loaded, warmup_done, device, last_inference_at, in_flight_requests }. 200 если status=ok, 503 если degraded/down. last_inference_at обновляется при каждом успешном /v1/embeddings.

## Inputs

- None specified.

## Expected Output

- `api/handlers/health.go`
- `api/handlers/health_test.go`

## Verification

Unit tests: status=ok 200, status=degraded 503, status=down 503. last_inference_at обновляется после inference. T-H-7 pass.
