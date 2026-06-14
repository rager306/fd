---
id: T02
parent: S04
milestone: M046-zqzcu6
key_files:
  - api/middleware/auth.go
  - api/middleware/auth_test.go
  - README.md
key_decisions:
  - Use fail-closed protected endpoint behavior rather than a shipped default API key or secret-bearing compose config.
duration: 
verification_result: passed
completed_at: 2026-06-14T18:59:37.819Z
blocker_discovered: false
---

# T02: Made protected endpoints fail closed when `FD_API_KEY` is missing and protected `/metrics`.

**Made protected endpoints fail closed when `FD_API_KEY` is missing and protected `/metrics`.**

## What Happened

Changed API key auth so public probe/docs endpoints and OPTIONS remain open, but protected endpoints reject requests when `FD_API_KEY` is empty instead of disabling authentication. Removed `/metrics` from the public auth exception list while keeping the path constant for tests and readability. Updated README controls table to document the fail-closed protected-endpoint behavior without committing or revealing any secret value.

## Verification

`cd api && go test ./middleware -run 'TestAPIKeyAuth|TestRateLimiterPrunesExpiredBuckets'` passed after implementation, including missing-key rejection and metrics protection tests.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./middleware -run 'TestAPIKeyAuth|TestRateLimiterPrunesExpiredBuckets'` | 0 | ✅ pass | 14000ms |

## Deviations

No opt-out auth-disable flag was added; the remediation intentionally changes the default protected-endpoint posture to fail closed.

## Known Issues

Runtime UAT remains in T04. Public probes remain unauthenticated by design.

## Files Created/Modified

- `api/middleware/auth.go`
- `api/middleware/auth_test.go`
- `README.md`
