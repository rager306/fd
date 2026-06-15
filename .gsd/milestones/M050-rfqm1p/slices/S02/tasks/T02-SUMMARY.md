---
id: T02
parent: S02
milestone: M050-rfqm1p
key_files:
  - tests/integration/api_test.go
  - benchmark-results/m050-s02-docker-e2e.md
key_decisions:
  - Keep integration secret separate from service runtime env var by using `FD_INTEGRATION_API_KEY`.
duration: 
verification_result: passed
completed_at: 2026-06-15T14:48:31.441Z
blocker_discovered: false
---

# T02: Расширен `tests/integration` до auth-aware Docker e2e suite.

**Расширен `tests/integration` до auth-aware Docker e2e suite.**

## What Happened

`tests/integration/api_test.go` теперь проверяет public `/live`, `/ready`, `/health`, auth fail-closed для embeddings/cache routes, authenticated `/metrics`, embeddings dimensions/batch, validation errors, missing-model compatibility, and cache HIT/flush/delete invalidation. Suite использует `FD_INTEGRATION_API_KEY` и не печатает секреты.

## Verification

`cd tests/integration && go test -v .` passed in no-key mode with 5 public/fail-closed checks; `cd api && go test ./...` remained passing with 295 tests.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd tests/integration && go test -v .` | 0 | ✅ pass | 7300ms |
| 2 | `cd api && go test ./...` | 0 | ✅ pass | 5400ms |

## Deviations

Metrics moved from public to authenticated checks after current runtime returned 401 without auth.

## Known Issues

Authenticated checks require `FD_INTEGRATION_API_KEY`; this is intentional.

## Files Created/Modified

- `tests/integration/api_test.go`
- `benchmark-results/m050-s02-docker-e2e.md`
