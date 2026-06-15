---
id: M049-7dn2gp
title: "Agent native observability and cache controls"
status: complete
completed_at: 2026-06-15T13:17:18.689Z
key_decisions:
  - Use input+dimensions for cache delete instead of non-reversible `:keyHash`.
  - Keep cache invalidation behind existing API key auth for solo deployment rather than add a separate admin-token architecture.
  - Use lightweight TEI /health and Redis PING dependency probes, not embedding inference, for /health dependency context.
  - Expose L1 cache occupancy in metrics and avoid Redis SCAN per scrape.
  - Defer AN-D trace hardening and broad AN-E/F config/options extraction until fd has a multi-tenant or externally managed policy need.
key_files:
  - api/cache/local.go
  - api/cache/redis.go
  - api/cache/tiered.go
  - api/handlers/cache.go
  - api/handlers/health.go
  - api/observability/metrics.go
  - api/main.go
  - benchmark-results/m049-s03-live-container-proof.md
  - benchmark-results/m049-issue-8-closure.md
lessons_learned:
  - For solo operation, explicit operator endpoints and cheap observability are more valuable than premature multi-tenant policy frameworks.
  - Runtime proof should verify cache behavior through actual cache headers, not just source tests.
---

# M049-7dn2gp: Agent native observability and cache controls

**M049 resolves the requested issue #8 solo-operator scope with cache invalidation controls plus health and metrics diagnostics, verified in a rebuilt container.**

## What Happened

M049 addressed the user-requested subset of GitHub issue #8 through GSD. S01 implemented AN-A by adding namespace-safe cache invalidation primitives and authenticated cache flush/delete routes. S02 implemented AN-B/AN-C by adding health last_error/dependency/capacity context and Prometheus runtime/cache gauges. S03 ran final gates, rebuilt the API container, proved the new behavior live against TEI and Redis, and wrote the issue #8 closure matrix. AN-D trace hardening was intentionally deferred per user direction. AN-E/F broad config/options extraction was avoided for solo deployment; only minimal seams needed for implemented health/cache surfaces were added. AN-G/H/I remain low-priority follow-ups outside this requested implementation scope.

## Success Criteria Results

- ✅ Issue #8 artifact saved: `documents/issue-8-current-m049.md`.
- ✅ AN-A implemented and live-verified: cache flush/delete primitives/routes, namespace-scoped Redis invalidation, authenticated live flush/delete cycles.
- ✅ AN-B/AN-C implemented and live-verified: `/health` dependency/capacity context, `/metrics` runtime/cache gauges.
- ✅ AN-D deferred and AN-E/F solo-scoped in D051 and closure matrix.
- ✅ Final tests/lint/govulncheck/UAT/container smoke/milestone validation passed.

## Definition of Done Results

- ✅ All slices complete: S01, S02, S03.
- ✅ Requirements validated: R040, R041, R042.
- ✅ Final `go test ./...`: 295 passed in 10 packages.
- ✅ Final golangci-lint: 0 issues.
- ✅ Final govulncheck: 0 reachable vulnerabilities.
- ✅ Docker Compose api/redis/tei healthy after rebuild.
- ✅ Live smoke: 5 passed, 0 failed.
- ✅ Milestone validation verdict: pass.

## Requirement Outcomes

| Requirement | Status | Proof |
|---|---|---|
| R040 | Validated | S01 implementation/tests and S03 live cache invalidation proof. |
| R041 | Validated | S02 implementation/tests and S03 live health/metrics proof. |
| R042 | Validated | D051 and closure matrix record solo-scope decision for AN-D/E/F. |

## Deviations

AN-G/H/I were not implemented because the user specifically requested AN-A and AN-B/C, with decisions for AN-D/E/F. Redis occupancy gauge was not added to avoid SCAN overhead; L1 occupancy is exposed instead.

## Follow-ups

Optional future cleanup: AN-G remove/wire `publicMetrics`, AN-H expose embedding headers via CORS, AN-I update README poll/single-tenant docs. Optional outward action after explicit confirmation: push commits and comment/close issue #8.
