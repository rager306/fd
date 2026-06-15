# S02: Runtime contract simplification — UAT

**Milestone:** M048-l4sctg
**Written:** 2026-06-15T11:19:39.873Z

# S02: Runtime contract simplification — UAT

**Milestone:** M048-l4sctg
**Written:** 2026-06-15

## UAT Type

- UAT mode: artifact-driven
- Why this mode is sufficient: S02 simplifies backend source contracts and health metadata. The observable outcome is source shape plus tests; no browser surface is involved.

## Preconditions

- `benchmark-results/m048-s02-runtime-contract-cleanup.md` exists.

## Smoke Test

Verify RuntimeHealth fields, shared embed interface, lifecycle state construction, and evidence artifact.

## Test Cases

### 1. RuntimeHealth active fields only

1. Inspect `api/handlers/health.go`.
2. **Expected:** inactive ONNX-only fields are absent from RuntimeHealth.

### 2. Shared embed interface

1. Inspect `api/embed/types.go`, `api/handlers/embeddings.go`, and `api/lifecycle/warmup.go`.
2. **Expected:** `embed.Embedder` exists and duplicate handler/lifecycle interface declarations are absent.

### 3. Explicit lifecycle state

1. Inspect `api/lifecycle/state.go` and `api/main.go`.
2. **Expected:** `defaultState`/`DefaultState()` are absent and main calls `lifecycle.NewState()`.

### 4. Evidence artifact complete

1. Inspect `benchmark-results/m048-s02-runtime-contract-cleanup.md`.
2. **Expected:** artifact covers #26, #29, #30, tests, and R038 validation.

## Edge Cases

- `/health.runtime` still includes active TEI fields: backend, model, dimensions, production_default, cache_namespace.
- Manual warmup handler still compiles against shared embedder interface.

## Failure Signals

- ONNX-only RuntimeHealth fields return.
- Duplicate interface declarations return.
- Lifecycle singleton returns.
- `go test ./...` fails.

## Requirements Proved By This UAT

- R038: active TEI runtime/lifecycle contracts are simplified and explicit.

## Not Proven By This UAT

- Validation error message and OpenAPI helper cleanup; these are S03.

## Notes for Tester

UAT evidence: `a242ec62-a82e-48e6-962f-f21f9b87bc28`, `7fe8e369-348d-4184-8292-880b5425da2b`, `346c2f5f-bc4b-4905-98ff-8f9edb6b37c4`, `7ffa9040-28ba-4f7f-aa9f-e50b58a036ed`.
