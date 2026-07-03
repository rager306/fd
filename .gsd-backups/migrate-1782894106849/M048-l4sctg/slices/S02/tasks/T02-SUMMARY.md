---
id: T02
parent: S02
milestone: M048-l4sctg
key_files:
  - api/embed/types.go
  - api/handlers/embeddings.go
  - api/handlers/batch.go
  - api/handlers/v1batch.go
  - api/handlers/batch_backend.go
  - api/lifecycle/warmup.go
  - api/lifecycle/state.go
  - api/main.go
  - api/handlers/health.go
  - api/handlers/health_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-15T11:17:03.760Z
blocker_discovered: false
---

# T02: Simplified runtime health, embed interface, and lifecycle state contracts.

**Simplified runtime health, embed interface, and lifecycle state contracts.**

## What Happened

Removed inactive ONNX-only fields from `RuntimeHealth` and deleted the stale ONNX runtime health test. Added shared `embed.Embedder` and updated handlers, batch helpers, lifecycle warmup, warmup handler, and integration tests to use it instead of duplicate `handlers.Embedder` and `lifecycle.WarmupModel` declarations. Removed lifecycle package default singleton and changed main to explicitly call `lifecycle.NewState()`.

## Verification

Focused `go test ./handlers ./lifecycle` passed with 101 tests. Full `go test ./...` passed with 280 tests. Static proof `d75568af-277e-40e2-a28b-e6ee373d28dd` passed for removed ONNX-only fields, duplicate interfaces, and lifecycle singleton.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./handlers ./lifecycle` | 0 | ✅ pass | 8700ms |
| 2 | `cd api && go test ./...` | 0 | ✅ pass | 29900ms |
| 3 | `gsd_exec d75568af-277e-40e2-a28b-e6ee373d28dd` | 0 | ✅ pass | 130ms |

## Deviations

Removed the stale ONNX health test instead of rewriting it because ONNX is intentionally not active in the current product/runtime path.

## Known Issues

S03 still needs validation message and OpenAPI helper cleanup.

## Files Created/Modified

- `api/embed/types.go`
- `api/handlers/embeddings.go`
- `api/handlers/batch.go`
- `api/handlers/v1batch.go`
- `api/handlers/batch_backend.go`
- `api/lifecycle/warmup.go`
- `api/lifecycle/state.go`
- `api/main.go`
- `api/handlers/health.go`
- `api/handlers/health_test.go`
