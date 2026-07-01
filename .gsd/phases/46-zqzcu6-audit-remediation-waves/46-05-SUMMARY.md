---
id: S05
parent: M046-zqzcu6
milestone: M046-zqzcu6
provides:
  - Issue #3 P1 #10 closed.
  - R031 validated.
  - S06 can focus on residual audit closure rather than LocalCache correctness.
requires:
  - slice: S01
    provides: Confirmed LocalCache P1 #10 finding and corrective direction.
affects:
  []
key_files:
  - api/cache/local.go
  - api/cache/local_test.go
  - api/main.go
  - benchmark-results/m046-s05-localcache-correctness.md
  - documents/issue-3-audit-remediation-plan-m046.md
  - .gsd/REQUIREMENTS.md
key_decisions:
  - Replace sync.Map plus separate counter with one mutex-owned map and derived length.
  - Add idempotent Close for LocalCache lifecycle instead of leaving background eviction goroutines unmanaged.
patterns_established:
  - Cache accounting should have one source of truth guarded by one lock; lifecycle-managed goroutines should expose idempotent Close.
observability_surfaces:
  - benchmark-results/m046-s05-localcache-correctness.md records red/green tests, race tests, gates, and static proof.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-14T19:23:30.804Z
blocker_discovered: false
---

# S05: LocalCache correctness

**Made LocalCache accounting and lifecycle deterministic with a mutex-owned map, derived size, idempotent Close, race-tested behavior, and API shutdown integration.**

## What Happened

S05 closed issue #3 P1 #10. The previous LocalCache used `sync.Map`, a separate mutex-protected `size` counter, and an eviction goroutine with no stop surface. The implementation now uses a single mutex-owned `map[string]l1Entry`, derives size from `len(data)`, refreshes overwrites in place, expires entries under the same lock, and enforces capacity while retaining the just-written key. `Close() error` is idempotent and stops the eviction loop via `sync.Once`, `stopCh`, and `doneCh`. API shutdown and Redis error paths now close the local cache. Targeted tests, race-enabled LocalCache tests, full gates, static proof, and artifact UAT passed. R031 was validated.

## Verification

Red tests first failed with `c.Close undefined`. After implementation, targeted LocalCache tests passed, `cd api && go test -race ./cache -run TestLocalCache` passed, `cd api && go test ./cache && go test ./...` passed with 44 cache tests and 281 total tests, golangci-lint reported 0 issues, govulncheck reported 0 reachable vulnerabilities, static proof `f124000a-5996-4c68-888d-1e31237c6d39` passed, completeness proof `a623b5a9-3cd3-4413-bb88-20ee70f64547` passed, and artifact-driven UAT PASS was saved.

## Requirements Advanced

None.

## Requirements Validated

- R031 — S05 tests, race tests, static proof, full gates, and artifact UAT prove LocalCache deterministic size/accounting/lifecycle behavior.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Used artifact-driven UAT rather than runtime HTTP UAT because S05 changes internal cache implementation and lifecycle, not external route behavior.

## Known Limitations

LocalCache still returns/stores byte slices without defensive copying, matching pre-existing behavior. S06 can decide whether that is in scope for residual hardening; it was not part of issue #3 P1 #10.

## Follow-ups

Proceed to S06 for residual P1 #6 and P2/P3 closure matrix, then milestone validation.

## Files Created/Modified

- `api/cache/local.go` — Refactored LocalCache to a mutex-owned map, derived size, idempotent Close, and stoppable eviction loop.
- `api/cache/local_test.go` — Added lifecycle and concurrent overwrite tests.
- `api/main.go` — Close LocalCache during shutdown and Redis error cleanup.
- `benchmark-results/m046-s05-localcache-correctness.md` — S05 evidence artifact.
- `documents/issue-3-audit-remediation-plan-m046.md` — Marked S05 done.
- `.gsd/REQUIREMENTS.md` — R031 validated.
