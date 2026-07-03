---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T03: Disabled trusted forwarded headers by default and bounded rate limiter key state.

Configure Gin to trust no forwarded proxies by default and add rate limiter cleanup so per-key state cannot grow unbounded indefinitely. Keep behavior configurable only if a future slice adds trusted proxy env policy; do not trust X-Forwarded-For by default.

## Inputs

- `api/main.go`
- `api/middleware/ratelimit.go`

## Expected Output

- `api/main.go`
- `api/middleware/ratelimit.go`
- `api/middleware/ratelimit_test.go`

## Verification

cd api && go test ./middleware -run TestRateLimiter && cd api && go test . -run 'TestRouterDoesNotTrustForwardedForByDefault'

## Observability Impact

Rate limiter state cleanup is covered by deterministic tests.
