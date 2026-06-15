---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Pinned health/metrics diagnostic gaps with red tests.

Add failing tests for last_error serialization, dependency blocks, in_flight_capacity, and Prometheus gauges for in-flight/cache occupancy.

## Inputs

- `api/handlers/health.go`
- `api/observability/metrics.go`
- `api/middleware/lifecycle.go`
- `documents/issue-8-current-m049.md`

## Expected Output

- `api/handlers/health_test.go`
- `api/observability/metrics_test.go`

## Verification

cd api && go test ./handlers ./observability (expected red before implementation).

## Observability Impact

Red tests pin diagnostic fields required by AN-B/AN-C.
