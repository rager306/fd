# S02: Runtime contract simplification

**Goal:** Resolve issue #7 findings #26, #29, and #30 by simplifying active TEI runtime health, unifying duplicate embedding interfaces, and removing the lifecycle default singleton.
**Demo:** Health and lifecycle contracts expose only active TEI runtime surfaces and one embedding interface contract.

## Must-Haves

- Static pre-fix evidence confirms ONNX-only RuntimeHealth fields, duplicate Embedder/WarmupModel declarations, and default lifecycle singleton exist.
- RuntimeHealth exposes active TEI fields only and health tests remain green.
- Embedding/warmup code uses one shared embedder interface without import cycles.
- `lifecycle.DefaultState` singleton is removed and main constructs state explicitly.
- `go test ./...` passes and R038 is validated.

## Proof Level

- This slice proves: Static symbol checks, focused health/lifecycle tests, full Go tests.

## Integration Closure

Same-host `/health.runtime` still exposes backend/model/dimensions/production_default/cache_namespace.

## Verification

- Health output becomes less misleading for future operators/agents.

## Tasks

- [x] **T01: Confirmed issue #7 runtime contract debt exists before S02 fixes.** `est:small`
  Run static checks proving issue #7 #26/#29/#30 are present before fixes: ONNX-only RuntimeHealth fields, duplicate interfaces, default singleton.
  - Files: `api/handlers/health.go`, `api/handlers/embeddings.go`, `api/lifecycle/warmup.go`, `api/lifecycle/state.go`
  - Verify: Static gsd_exec check should PASS for pre-fix presence.

- [x] **T02: Simplified runtime health, embed interface, and lifecycle state contracts.** `est:medium`
  Remove ONNX-only RuntimeHealth fields and stale ONNX health tests/fields, define one shared embed interface in `api/embed`, update handlers/lifecycle to use it, and remove lifecycle default singleton in favor of explicit `NewState` in main.
  - Files: `api/embed/types.go`, `api/handlers/embeddings.go`, `api/lifecycle/warmup.go`, `api/lifecycle/state.go`, `api/main.go`, `api/handlers/health.go`, `api/handlers/health_test.go`, `api/main_test.go`
  - Verify: cd api && go test ./handlers ./lifecycle ./...

- [x] **T03: Recorded S02 runtime contract evidence and validated R038.** `est:small`
  Write S02 evidence artifact, validate R038, run full tests/static post-check, and complete S02.
  - Files: `benchmark-results/m048-s02-runtime-contract-cleanup.md`, `.gsd/REQUIREMENTS.md`
  - Verify: cd api && go test ./... plus static post-cleanup check.

## Files Likely Touched

- api/handlers/health.go
- api/handlers/embeddings.go
- api/lifecycle/warmup.go
- api/lifecycle/state.go
- api/embed/types.go
- api/main.go
- api/handlers/health_test.go
- api/main_test.go
- benchmark-results/m048-s02-runtime-contract-cleanup.md
- .gsd/REQUIREMENTS.md
