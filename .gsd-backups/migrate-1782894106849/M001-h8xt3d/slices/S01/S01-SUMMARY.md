---
id: S01
parent: M001-h8xt3d
milestone: M001-h8xt3d
provides:
  - Safer cache behavior for API validation and handler test work in S02.
requires:
  []
affects:
  - S02
  - S03
  - S04
key_files:
  - api/cache/tiered.go
  - api/cache/redis.go
  - api/cache/redis_binary_test.go
  - api/cache/tiered_cache_test.go
key_decisions:
  - Keep Redis write failures best-effort in TieredCache after successful validation, preserving cache-aside behavior while preventing serialization panics.
patterns_established:
  - Use dimension-aware keys for any cache layer that can store multiple vector dimensions for the same text.
  - Validate embedding length at serialization boundaries rather than trusting upstream model output.
observability_surfaces:
  - Short-vector failures now return explicit errors containing observed and requested dimensions.
drill_down_paths:
  - .gsd/milestones/M001-h8xt3d/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M001-h8xt3d/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M001-h8xt3d/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T06:52:39.186Z
blocker_discovered: false
---

# S01: Cache correctness and panic safety

**S01 fixed cache dimension isolation and short-vector panic safety.**

## What Happened

S01 addressed the highest-risk cache correctness findings. TieredCache now creates a dimension-aware local key for L1 and singleflight, preventing same-text 512d and 1024d calls from sharing cached or in-flight results. marshalEmbedding now validates dimensions and vector length, returning errors instead of panicking when the embedding is too short. Tests were updated and expanded to prove both the binary serialization error behavior and TieredCache dimension isolation.

## Verification

Targeted cache tests passed with 13 tests; full short suite passed with 38 tests across 4 packages.

## Requirements Advanced

- Review remediation high-risk cache findings advanced with tests. — 

## Requirements Validated

- Cache no longer shares same-text 512d and 1024d local/singleflight results — proved by TestTieredCache_GetOrLoad_SeparatesDimensionsForSameText.
- Short embeddings now return errors instead of panicking — proved by TestMarshalEmbedding_ShortVectorReturnsError and TestTieredCache_GetOrLoad_ReturnsErrorForShortEmbedding.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

GitNexus/LSP impact analysis was unavailable for this Go repo, so text-search blast radius was used and documented.

## Known Limitations

Redis SetBytes command errors remain best-effort in TieredCache to avoid failing successful TEI inference due to cache write issues.

## Follow-ups

Continue with S02 API validation and handler tests.

## Files Created/Modified

- `api/cache/tiered.go` — Added dimension-aware L1/singleflight keys and safe marshaling propagation.
- `api/cache/redis.go` — Changed marshalEmbedding to return validation errors and fixed SetBytes dimension-aware keying.
- `api/cache/redis_binary_test.go` — Updated binary serialization tests for error-returning marshal function and added short-vector tests.
- `api/cache/tiered_cache_test.go` — Added TieredCache tests for 512d/1024d same-text isolation and short-vector error behavior.
