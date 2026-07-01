---
id: S01
parent: M049-7dn2gp
milestone: M049-7dn2gp
provides:
  - AN-A implemented at source/test level.
  - R040 advanced pending live runtime proof.
requires:
  []
affects:
  []
key_files:
  - api/cache/local.go
  - api/cache/redis.go
  - api/cache/tiered.go
  - api/handlers/cache.go
  - api/main.go
  - benchmark-results/m049-s01-cache-invalidation.md
key_decisions:
  - Use input+dimensions for delete route instead of non-reversible key hashes.
  - Keep cache invalidation behind existing API key auth rather than introduce a separate admin token in solo scope.
patterns_established:
  - Destructive cache actions must be namespace-scoped and auth-gated.
  - Runtime cache proof should be delayed until aggregate container verification to avoid rebuilding for every slice.
observability_surfaces:
  - Authenticated cache flush/delete routes.
  - S01 evidence artifact documenting cache invalidation behavior.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-15T12:54:56.953Z
blocker_discovered: false
---

# S01: Cache invalidation controls

**fd now has tested cache invalidation primitives and authenticated cache flush/delete routes for solo operator use.**

## What Happened

S01 implemented issue #8 AN-A. The active cache stack now has first-class invalidation primitives: LocalCache can flush and report size; RedisCache can delete single dimensioned entries and flush only its namespace via SCAN/DEL; TieredCache delegates delete/flush across L1 and L2. The HTTP surface adds `POST /v1/cache/flush` and `POST /v1/cache/delete`, wired behind the existing global API key middleware. The delete route accepts string or string-array input plus dimensions, defaulting to 1024 and validating the existing 512/1024 dimensions. We deliberately avoided `:keyHash` deletion because the current short hash is not a reversible or unique operator input; input+dimensions maps directly to the real cache key derivation.

## Verification

Red tests first failed on missing primitives/handler. Green verification passed: `cd api && go test ./cache ./handlers` passed with 127 tests; `cd api && go test ./...` passed with 293 tests. Static proof `3670b28f-8bce-433e-8306-987102db98cb` verified source invariants. S01 UAT passed with evidence `94ea4377-4e0a-4327-a167-76d5bcf0404c`, `6d55b34d-006b-431c-8045-cb8e5f639981`, and `8707655c-51b1-452e-9af3-1efd9ba08dda`.

## Requirements Advanced

- R040 — Implemented and tested cache invalidation primitives/routes; runtime proof remains in S03.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Runtime cache HIT->flush->MISS proof deferred to S03; R040 remains active until then.

## Known Limitations

No per-key hash route was added because the current short hash is not a safe deletion identifier. Delete by input+dimensions is the implemented solo-operator action.

## Follow-ups

S03 must rebuild the container and prove authenticated HIT->flush->MISS behavior live.

## Files Created/Modified

- `api/cache/local.go` — Added Flush and Size.
- `api/cache/redis.go` — Added namespace-scoped Delete and FlushNamespace.
- `api/cache/tiered.go` — Added Delete, Flush, and LocalSize.
- `api/handlers/cache.go` — Added cache invalidation HTTP handler.
- `api/main.go` — Registered cache invalidation routes.
