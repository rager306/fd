# M046 Issue #3 Audit Remediation Plan

Updated: 2026-06-14

## Purpose

Issue #3 reports 32 findings across security, API contract, performance, reliability, and cache correctness. M046 handles this in waves: first validate the findings against current code, then fix the highest-risk confirmed clusters with targeted tests.

## Confirmed Root Decisions

### 1. Batch endpoints evolved outside the hardened `/v1/embeddings` boundary

`/v1/embeddings` has a clear middleware boundary: validation, rate limiting, lifecycle/capacity, and OpenAI-style errors before backend work. `/v1/batch` and `/embeddings/batch` were added as parallel endpoint surfaces and did not inherit the full boundary.

Impact:
- P0 #2 and #3 are confirmed.
- P1 #4 and #5 are confirmed because backend work is still one TEI call per input on cache misses.

Corrective direction:
- S02 adds request/abuse guardrails before backend work.
- S03 reshapes backend calls to collect misses and embed in bounded chunks.

### 2. Same-host local assumptions leaked into default exposure behavior

The service contract targets same-host consumers, but compose publishes the API on all interfaces and empty `FD_API_KEY` disables auth. PR #2 correctly made health/readiness endpoints public for probes, but inference/admin/diagnostic exposure remains a policy risk.

Impact:
- P0 #1 is confirmed as a policy risk outside strict loopback-only deployments.
- P1 #7 and #8 are confirmed or likely valid depending on the selected exposure stance.

Corrective direction:
- S04 must define the default exposure policy before changing code.
- Candidate safe baseline: bind local compose to loopback by default and require explicit opt-in or key for non-loopback exposure; keep `/live`, `/ready`, `/health`, and `/v1/healthcheck` public.
- Metrics/traces should be protected unless explicitly enabled for local diagnostics.

### 3. LocalCache prioritized simple concurrency over deterministic lifecycle and accounting

`LocalCache` uses `sync.Map`, a separate mutex-protected size counter, and an always-running eviction goroutine. This is simple and fast enough for happy paths, but hard to reason about under concurrent overwrite/delete/expiry and has no shutdown method.

Impact:
- P1 #10 is confirmed as a correctness/lifecycle risk.

Corrective direction:
- S05 should either harden LocalCache with a single lock-owned map and `Close`, or revive a true bounded LRU implementation.
- Verification must include race-enabled cache tests and capacity tests.

## Wave Plan

| Wave | Slice | Fixes | Output |
|---|---|---|---|
| 1 | S01 Audit validation map | Confirms P0/P1 and root decisions | This document plus `benchmark-results/m046-s01-audit-validation.md` |
| 2 | S02 Batch endpoint guardrails | P0 #2, P0 #3 | Tests proving invalid/oversized batch requests are rejected before backend work |
| 3 | S03 Batch backend work shaping | P1 #4, P1 #5, maybe P1 #6 | Tests proving bounded TEI call counts and preserved response order |
| 4 | S04 Exposure posture policy | P0 #1, P1 #7, P1 #8 and related exposure items | Auth/compose policy tests and docs |
| 5 | S05 LocalCache correctness | P1 #10 | Race/capacity/lifecycle tests |
| 6 | S06 Audit closure | Remaining P1 plus P2/P3 triage | Closure matrix and future requirements/non-goals |

## S02 Starting Checklist

Files:
- `api/main.go`
- `api/handlers/batch.go`
- `api/handlers/v1batch.go`
- `api/middleware/validation.go`
- `api/middleware/ratelimit.go`
- `api/middleware/lifecycle.go`
- `api/handlers/*batch*_test.go`
- `tests/integration/api_test.go` if HTTP-level proof is needed

Required behavior to prove before implementation:
- `/v1/batch` rejects oversized bodies before handler/backend work.
- `/v1/batch` rejects too many strings and too-long strings before backend work.
- `/embeddings/batch` rejects oversized bodies before handler/backend work.
- `/embeddings/batch` rejects too many strings and too-long strings before backend work.
- lifecycle/capacity and rate-limit middleware apply to both batch endpoints where appropriate.
- valid batch requests still work.

## Policy Decisions To Revisit In S04

- Should local compose bind `api` to `127.0.0.1:8000` by default?
- Should `GIN_MODE=release` fail startup when `FD_API_KEY` is empty unless an explicit `FD_ALLOW_UNAUTHENTICATED_LOCAL=true` style opt-in is present?
- Should `/metrics` and `/v1/traces` require auth by default, with a local diagnostics opt-in?

No GitHub comments or issue status changes should be made without explicit user confirmation.
