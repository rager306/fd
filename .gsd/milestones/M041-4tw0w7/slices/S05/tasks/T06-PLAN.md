---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T06: Added `/v1/traces` in-memory request trace ring buffer and endpoint.

api/observability/traces.go: in-memory ring buffer (последние 100 requests) с timestamp, latency, status, model_id, request_id, path, dimensions. GET /v1/traces возвращает JSON массив. Использует request_id из headers middleware (S03). Опционально через FD_TRACES_ENABLED=true (default true).

## Inputs

- None specified.

## Expected Output

- `api/observability/traces.go`
- `api/observability/traces_test.go`

## Verification

Unit tests: после 5 requests GET /v1/traces → 200 с 5 entries. Каждая entry содержит timestamp, latency_ms, status, model_id, request_id, path.
