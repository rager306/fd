---
id: T01
parent: S04
milestone: M046-zqzcu6
key_files:
  - api/middleware/auth_test.go
  - api/middleware/ratelimit_test.go
  - api/main_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T18:57:15.825Z
blocker_discovered: false
---

# T01: Added red tests for S04 exposure posture contracts.

**Added red tests for S04 exposure posture contracts.**

## What Happened

Added tests that define the intended S04 security posture: missing API key rejects protected endpoints, `/metrics` is protected, only probe/docs endpoints remain public, rate limiter state is pruned instead of growing unbounded, and the router does not trust forwarded headers by default. The first targeted test run is red because `maxRateLimitKeys` does not exist yet; the main proxy test is red because `configureTrustedProxies` does not exist yet.

## Verification

`cd api && go test ./middleware -run 'TestAPIKeyAuth|TestRateLimiterPrunesExpiredBuckets'` failed with undefined `maxRateLimitKeys`; `cd api && go test . -run 'TestRouterDoesNotTrustForwardedForByDefault'` failed with undefined `configureTrustedProxies`. These are expected red failures before implementation.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./middleware -run 'TestAPIKeyAuth|TestRateLimiterPrunesExpiredBuckets'` | 1 | ✅ expected red fail | 1000ms |
| 2 | `cd api && go test . -run 'TestRouterDoesNotTrustForwardedForByDefault'` | 1 | ✅ expected red fail | 1000ms |

## Deviations

Red state is compile-red for the new limiter/proxy seams, while auth behavior will be exercised once the package compiles.

## Known Issues

Implementation pending in T02/T03.

## Files Created/Modified

- `api/middleware/auth_test.go`
- `api/middleware/ratelimit_test.go`
- `api/main_test.go`
