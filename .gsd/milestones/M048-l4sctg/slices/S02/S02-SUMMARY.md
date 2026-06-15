---
id: S02
parent: M048-l4sctg
milestone: M048-l4sctg
provides:
  - R038 validated.
  - Issue #7 findings #26, #29, and #30 closed.
requires:
  - slice: S01
    provides: Cache cleanup baseline and issue input artifact.
affects:
  []
key_files:
  - api/embed/types.go
  - api/handlers/health.go
  - api/handlers/health_test.go
  - api/handlers/embeddings.go
  - api/handlers/batch.go
  - api/handlers/v1batch.go
  - api/handlers/batch_backend.go
  - api/handlers/warmup.go
  - api/lifecycle/warmup.go
  - api/lifecycle/state.go
  - api/main.go
  - benchmark-results/m048-s02-runtime-contract-cleanup.md
key_decisions:
  - Centralize the inference interface in package embed, where the concrete TEI client lives.
  - Remove ONNX-only health fields entirely because TEI-only is the active product path.
patterns_established:
  - Runtime health types should describe active product behavior, not dormant research branches.
observability_surfaces:
  - Health endpoint now exposes fewer misleading fields; S02 evidence artifact records proof.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-15T11:19:39.873Z
blocker_discovered: false
---

# S02: Runtime contract simplification

**Runtime health, embedding interface, and lifecycle state contracts now reflect the active TEI-only path.**

## What Happened

S02 resolved issue #7 findings #26, #29, and #30. Pre-fix proof confirmed stale ONNX-only RuntimeHealth fields, duplicate Embedder/WarmupModel interfaces, and lifecycle default singleton. The fix removed inactive ONNX health fields and stale ONNX health tests, introduced shared `embed.Embedder`, updated handlers/lifecycle/warmup/main/tests to consume the shared interface, removed `defaultState`/`DefaultState`, and made main explicitly construct lifecycle state with `lifecycle.NewState()`. R038 was validated.

## Verification

Pre-fix proof `5ef6afae-43e0-41bc-94a3-dd43253cec50` passed. Focused `go test ./handlers ./lifecycle` passed with 101 tests. Full `go test ./...` passed with 280 tests. Post-cleanup static proof `d75568af-277e-40e2-a28b-e6ee373d28dd` passed. Artifact completeness `2b24ced0-db5a-48c0-8fdd-6e3945e5c037` passed. UAT PASS saved with evidence `a242ec62-a82e-48e6-962f-f21f9b87bc28`, `7fe8e369-348d-4184-8292-880b5425da2b`, `346c2f5f-bc4b-4905-98ff-8f9edb6b37c4`, and `7ffa9040-28ba-4f7f-aa9f-e50b58a036ed`.

## Requirements Advanced

None.

## Requirements Validated

- R038 — S02 tests and static proof validate RuntimeHealth simplification, shared embed interface, and explicit lifecycle state construction.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Removed the stale ONNX health test rather than rewriting it because ONNX is not active in the current product/runtime path.

## Known Limitations

S03 issue #7 findings remain open.

## Follow-ups

Proceed to S03 validation/OpenAPI polish and milestone closure.

## Files Created/Modified

- `api/embed/types.go` — Added shared Embedder interface.
- `api/handlers/health.go` — Removed inactive ONNX-only RuntimeHealth fields.
- `api/lifecycle/state.go` — Removed default singleton.
- `api/main.go` — Uses explicit lifecycle.NewState and shared embed interface.
- `benchmark-results/m048-s02-runtime-contract-cleanup.md` — S02 evidence artifact.
