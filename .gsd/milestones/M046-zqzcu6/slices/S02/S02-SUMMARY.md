---
id: S02
parent: M046-zqzcu6
milestone: M046-zqzcu6
provides:
  - Validated R029 batch endpoint guardrails.
  - P0 #2 and #3 closed for issue #3.
  - S03 can now optimize batch backend calls within bounded inputs.
requires:
  - slice: S01
    provides: Verified P0/P1 issue #3 inventory and S02 starting checklist.
affects:
  []
key_files:
  - api/main.go
  - api/middleware/validation.go
  - api/handlers/batch_limits.go
  - api/handlers/batch.go
  - api/handlers/v1batch.go
  - api/handlers/embeddings_integration_test.go
  - api/handlers/v1batch_test.go
  - api/middleware/validation_test.go
  - benchmark-results/m046-s02-batch-guardrails.md
  - documents/issue-3-audit-remediation-plan-m046.md
key_decisions:
  - Use route-level body cap plus handler-specific batch validation instead of reusing `/v1/embeddings` middleware on incompatible batch JSON shapes.
patterns_established:
  - Batch endpoints with distinct JSON shapes should share generic request body and lifecycle guardrails while validating their own shape before backend work.
  - DoS-class findings should be proven by tests/static probes rather than abusive runtime load.
observability_surfaces:
  - benchmark-results/m046-s02-batch-guardrails.md records red test, implementation proof, route guardrail check, and gates.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-14T16:11:20.670Z
blocker_discovered: false
---

# S02: Batch endpoint guardrails

**Hardened `/v1/batch` and `/embeddings/batch` with body, rate-limit, lifecycle, and input guardrails before backend work.**

## What Happened

S02 fixed the confirmed issue #3 P0 batch endpoint defects. The implementation added a shared request body cap middleware for non-embeddings JSON shapes, mounted both batch routes with body cap, user rate limiting, and lifecycle/capacity middleware, and added handler-specific validation for legacy `inputs` and v1 nested `batches` shapes before cache or TEI work. The legacy batch endpoint now enforces non-empty inputs, max 128 inputs, and max 2048 chars per input. `/v1/batch` retains its group and inner-batch limits and now also rejects too-long nested strings. Both handlers surface `http.MaxBytesError` as `payload_too_large`. R029 was validated. P1 N+1 backend call shaping remains intentionally scoped to S03 after inputs are bounded.

## Verification

Red tests failed before implementation with undefined `maxBatchInputChars`, then targeted handler/middleware tests passed. `cd api && go test ./...` passed across 9 packages; golangci-lint v2.12.2 reported 0 issues; govulncheck v1.3.0 reported 0 reachable vulnerabilities. Static route guardrail check `070f8b99-679d-45e2-a90e-597c350f6837` passed. Runtime UAT after rebuilding `api` verified `/ready` and `/health`, `/v1/batch` 413 `input_too_long`, `/embeddings/batch` 413 `input_too_long`, and valid batch smokes for both endpoints.

## Requirements Advanced

None.

## Requirements Validated

- R029 — S02 tests, route guardrail check, and runtime UAT prove both batch endpoints enforce body, rate-limit, lifecycle, and input guardrails before backend work.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S02 did not attempt to solve N+1 batch TEI calls. That remains S03 by design, after S02 bounds input/work surfaces.

## Known Limitations

Rate-limiter trusted proxy and diagnostics exposure concerns remain S04/S06 scope. Batch endpoints are now guarded, but backend calls are still per-input on cache misses until S03.

## Follow-ups

Proceed to S03 to shape batch backend work and reduce N+1 TEI calls.

## Files Created/Modified

- `api/main.go` — Mounted both batch routes with body cap, rate-limit, and lifecycle/capacity guardrails.
- `api/middleware/validation.go` — Added `LimitRequestBody` shared body cap middleware.
- `api/handlers/batch_limits.go` — Added shared batch input limits.
- `api/handlers/batch.go` — Added legacy batch validation and MaxBytesError handling before backend work.
- `api/handlers/v1batch.go` — Added nested input length validation and MaxBytesError handling.
- `api/handlers/embeddings_integration_test.go` — Added legacy batch too-long-input rejection test.
- `api/handlers/v1batch_test.go` — Added v1 batch too-long-input rejection test.
- `api/middleware/validation_test.go` — Added request body cap middleware test.
- `benchmark-results/m046-s02-batch-guardrails.md` — S02 verification evidence.
- `documents/issue-3-audit-remediation-plan-m046.md` — Marked S02 wave complete and left N+1 for S03.
