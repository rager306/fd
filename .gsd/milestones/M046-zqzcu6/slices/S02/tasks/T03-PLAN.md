---
estimated_steps: 1
estimated_files: 4
skills_used: []
---

# T03: Mounted both batch routes with body, rate-limit, and lifecycle guardrails.

Update route registration so `/v1/batch` and `/embeddings/batch` pass through appropriate validation, user rate limit, and lifecycle/capacity gates before handlers. Preserve PR #2 public health behavior.

## Inputs

- `api/main.go`
- `api/middleware/ratelimit.go`
- `api/middleware/lifecycle.go`

## Expected Output

- `api/main.go`
- `api/main_test.go`

## Verification

Route-level tests or integration tests show batch endpoints reject while service is not ready/over capacity and rate-limit middleware is applied when enabled.

## Observability Impact

Middleware ordering remains explicit and documented in code comments where needed.
