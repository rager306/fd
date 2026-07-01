---
id: S03
parent: M046-zqzcu6
milestone: M046-zqzcu6
provides:
  - Issue #3 P1 #4 and #5 closed.
  - R029 extended to include bounded backend work proof.
  - S04 can focus on auth/exposure without carrying batch DoS/N+1 risk.
requires:
  - slice: S02
    provides: Batch endpoints are guarded and inputs are bounded before backend chunking.
affects:
  []
key_files:
  - api/handlers/batch_backend.go
  - api/handlers/batch.go
  - api/handlers/v1batch.go
  - api/handlers/v1batch_test.go
  - api/handlers/embeddings_integration_test.go
  - benchmark-results/m046-s03-batch-backend-chunking.md
  - documents/issue-3-audit-remediation-plan-m046.md
key_decisions:
  - Keep S03 package-local to handlers and do not change the cache interface; S05 will handle LocalCache internals separately.
patterns_established:
  - Batch cache-miss shaping: collect misses with `GetIfPresent`, call embedder once per bounded miss chunk, backfill with `Set`, then restore original response order.
observability_surfaces:
  - benchmark-results/m046-s03-batch-backend-chunking.md records red/green call-count proof, static proof, and quality gates.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-14T16:26:41.630Z
blocker_discovered: false
---

# S03: Batch backend work shaping

**Reduced batch endpoint TEI work amplification by batching cache misses and preserving cache/order behavior.**

## What Happened

S03 fixed issue #3 P1 #4 and #5. The legacy `/embeddings/batch` endpoint and OpenAI-compatible `/v1/batch` endpoint no longer call TEI once per input on cache misses. Both now use package-local miss chunking: cache hits are read with `GetIfPresent`, misses are embedded in bounded chunks, results are backfilled with `Set`, and output order is restored. Legacy `/embeddings/batch` keeps its base64-by-default response shape; `/v1/batch` preserves nested response shape. Red tests first proved the problem with 4 calls instead of 1 for legacy batch and 8 calls instead of 2 for v1 batch. After implementation, focused tests also prove repeated identical requests are served from cache without extra embedder calls. Runtime UAT after rebuilding the API verified valid batch behavior and S02 rejection regression safety.

## Verification

Focused call-count tests passed after implementation; `cd api && go test ./handlers` passed; `cd api && go test ./...` passed across all packages; golangci-lint v2.12.2 reported 0 issues; govulncheck v1.3.0 reported 0 reachable vulnerabilities; static proof `6591611c-d4d4-4485-b17e-ac2be3aa5d6d` confirmed both batch handlers use miss chunking and no longer use `GetOrLoad`; runtime UAT PASS was saved with evidence IDs `cb8a0f47-c9f4-4daa-ba02-b68152bb85ac`, `db1bbf65-af81-41ea-b67c-bcb1f74c6efc`, `e48a4774-94ed-4d02-9ab2-cffba624437e`, and `ba43c8b5-9581-4d65-99ee-6d00f568b4b0`.

## Requirements Advanced

None.

## Requirements Validated

- R029 — S02 validated batch guardrails; S03 validated bounded backend work through call-count tests, static proof, full gates, and runtime UAT.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S03 did not address P1 #6. That finding concerns `/v1/embeddings` Redis cache-peek sequencing and is left for later triage/closure rather than conflated with batch endpoint N+1 work.

## Known Limitations

S04/S05/S06 remain pending. Runtime UAT does not directly count TEI calls in production; call-count reduction is proven by unit tests and static proof.

## Follow-ups

Proceed to S04 exposure/auth posture policy for P0 #1 plus P1 #7/#8.

## Files Created/Modified

- `api/handlers/batch_backend.go` — Added shared batch cache-miss chunking helper.
- `api/handlers/batch.go` — Changed legacy batch endpoint to use bounded miss chunking and preserve legacy response encoding.
- `api/handlers/v1batch.go` — Changed v1 batch endpoint to use one embedder call per bounded inner batch/miss group.
- `api/handlers/v1batch_test.go` — Added call-count/cache-hit tests for v1 batch endpoint.
- `api/handlers/embeddings_integration_test.go` — Added call-count/cache-hit tests for legacy batch endpoint.
- `benchmark-results/m046-s03-batch-backend-chunking.md` — S03 evidence artifact.
- `documents/issue-3-audit-remediation-plan-m046.md` — Marked S03 done for P1 #4/#5 and left P1 #6 for later triage.
