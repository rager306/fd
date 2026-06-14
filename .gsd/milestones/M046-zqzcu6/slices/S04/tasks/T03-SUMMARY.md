---
id: T03
parent: S04
milestone: M046-zqzcu6
key_files:
  - api/main.go
  - api/main_test.go
  - api/middleware/ratelimit.go
  - api/middleware/ratelimit_test.go
key_decisions:
  - Default Gin trusted proxies to nil so rate limiting uses the direct peer address unless a future explicit trusted-proxy policy is implemented.
duration: 
verification_result: passed
completed_at: 2026-06-14T18:59:52.593Z
blocker_discovered: false
---

# T03: Disabled trusted forwarded headers by default and bounded rate limiter key state.

**Disabled trusted forwarded headers by default and bounded rate limiter key state.**

## What Happened

Added `configureTrustedProxies` during router setup so Gin does not trust spoofable `X-Forwarded-For` headers by default. Added `maxRateLimitKeys` plus expired-bucket pruning and oldest-bucket eviction before inserting new limiter keys so the in-memory token-bucket map cannot grow unbounded indefinitely. Tests verify both the direct peer IP behavior and limiter pruning behavior.

## Verification

`cd api && go test ./middleware -run 'TestAPIKeyAuth|TestRateLimiterPrunesExpiredBuckets'` passed and `cd api && go test . -run 'TestRouterDoesNotTrustForwardedForByDefault'` passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./middleware -run 'TestAPIKeyAuth|TestRateLimiterPrunesExpiredBuckets'` | 0 | ✅ pass | 14000ms |
| 2 | `cd api && go test . -run 'TestRouterDoesNotTrustForwardedForByDefault'` | 0 | ✅ pass | 11000ms |

## Deviations

No trusted proxy env configuration was added in S04; defaulting to no trusted forwarded headers is the safe baseline. A future slice can add explicit trusted proxy configuration if deployment requires it.

## Known Issues

None for T03.

## Files Created/Modified

- `api/main.go`
- `api/main_test.go`
- `api/middleware/ratelimit.go`
- `api/middleware/ratelimit_test.go`
