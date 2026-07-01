---
id: S04
parent: M046-zqzcu6
milestone: M046-zqzcu6
provides:
  - Issue #3 P0 #1 closed.
  - Issue #3 P1 #7 closed.
  - Issue #3 P1 #8 closed.
  - R030 validated.
requires:
  - slice: S01
    provides: Confirmed exposure posture findings and public health/readiness decision.
affects:
  []
key_files:
  - api/middleware/auth.go
  - api/middleware/auth_test.go
  - api/middleware/ratelimit.go
  - api/middleware/ratelimit_test.go
  - api/main.go
  - api/main_test.go
  - README.md
  - benchmark-results/m046-s04-exposure-posture.md
  - documents/issue-3-audit-remediation-plan-m046.md
key_decisions:
  - Fail closed for protected endpoints when `FD_API_KEY` is missing instead of shipping a default key or secret-bearing compose config.
  - Default Gin trusted proxies to nil; future explicit proxy policy can be added separately if deployment requires it.
patterns_established:
  - Exposure-sensitive endpoints should have explicit public carve-outs; absence of a secret must not disable auth.
  - Client identity for rate limiting should not trust forwarded headers without an explicit trusted-proxy policy.
observability_surfaces:
  - benchmark-results/m046-s04-exposure-posture.md records red/green tests, static proof, gates, and runtime UAT.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-14T19:08:50.671Z
blocker_discovered: false
---

# S04: Exposure posture policy

**Protected endpoints now fail closed without `FD_API_KEY`, `/metrics` is protected, forwarded headers are not trusted by default, and rate limiter state is bounded.**

## What Happened

S04 closed the issue #3 exposure posture wave. Auth now keeps only probes/docs/OpenAPI public and rejects protected endpoints when `FD_API_KEY` is missing instead of treating an empty key as auth disabled. `/metrics` was removed from the auth public path set. Gin router setup now calls `SetTrustedProxies(nil)` so `ClientIP()` does not trust spoofable `X-Forwarded-For` headers by default. The in-memory rate limiter now has `maxRateLimitKeys` plus expired-bucket pruning and oldest-bucket eviction before adding new keys, preventing indefinite per-key map growth. README documents the fail-closed protected-endpoint posture without adding any secret value. R030 was validated.

## Verification

Red tests first failed on missing `maxRateLimitKeys` and `configureTrustedProxies`. After implementation, targeted middleware and main tests passed. `cd api && go test ./...` passed with 279 tests; golangci-lint reported 0 issues; govulncheck reported 0 reachable vulnerabilities; static proof `15cb6196-085b-4882-8e4b-18d17008ee4d` passed; runtime UAT after API rebuild passed for public probes, protected inference 401, protected metrics 401, and public OpenAPI.

## Requirements Advanced

None.

## Requirements Validated

- R030 — S04 tests, static proof, gates, and runtime UAT prove protected endpoints fail closed, metrics is protected, forwarded headers are not trusted by default, and limiter state is bounded.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Did not add trusted proxy env configuration in S04; the secure default is no trusted forwarded headers. Did not runtime-test authorized inference because no secret should be introduced or echoed in artifacts; unit tests cover correct bearer token behavior.

## Known Limitations

Local compose with no `FD_API_KEY` now intentionally returns 401 for protected endpoints. Users must configure `FD_API_KEY` outside git to use inference endpoints. S05/S06 remain pending.

## Follow-ups

Proceed to S05 for LocalCache correctness. S06 should triage residual P1 #6 and remaining P2/P3 findings.

## Files Created/Modified

- `api/middleware/auth.go` — Fail-closed auth behavior and removal of metrics from public auth paths.
- `api/middleware/auth_test.go` — Tests for missing key rejection, metrics protection, public carve-outs, and valid bearer token behavior.
- `api/middleware/ratelimit.go` — Bounded limiter key state with pruning/oldest eviction.
- `api/middleware/ratelimit_test.go` — Limiter pruning/cap test.
- `api/main.go` — Router trusted-proxy configuration defaults to nil.
- `api/main_test.go` — Test proving forwarded headers are not trusted by default.
- `README.md` — Updated FD_API_KEY operational posture.
- `benchmark-results/m046-s04-exposure-posture.md` — S04 evidence artifact.
- `documents/issue-3-audit-remediation-plan-m046.md` — Marked S04 done for P0 #1 and P1 #7/#8.
