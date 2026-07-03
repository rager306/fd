---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Added red tests for S04 exposure posture contracts.

Add failing tests proving empty `FD_API_KEY` must reject protected endpoints, public probes remain open, `/metrics` requires auth, forwarded headers are not trusted by default, and rate limiter state is pruned/bounded.

## Inputs

- `api/middleware/auth.go`
- `api/middleware/ratelimit.go`
- `api/main.go`

## Expected Output

- `api/middleware/auth_test.go`
- `api/middleware/ratelimit_test.go`
- `api/main_test.go`

## Verification

cd api && go test ./middleware -run 'TestAPIKeyAuth|TestRateLimiter' && cd api && go test . -run 'TestRouterDoesNotTrustForwardedForByDefault'

## Observability Impact

Tests become executable proof for issue #3 P0 #1 and P1 #7/#8.
