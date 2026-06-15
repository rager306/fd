# M048 S02 Runtime Contract Cleanup Evidence

Captured: 2026-06-15

## Scope

S02 covers GitHub issue #7 findings:

- #26 `RuntimeHealth` carries ONNX-only fields that are never populated in the active TEI runtime.
- #29 `Embedder` and `WarmupModel` are the same interface under two names.
- #30 `lifecycle.defaultState` singleton is redundant.

## Pre-fix Evidence

Static proof:

```text
gsd_exec 5ef6afae-43e0-41bc-94a3-dd43253cec50
PASS pre-fix issue #7 #26/#29/#30 runtime contract debt present
```

Confirmed:

- ONNX-only RuntimeHealth fields existed.
- `handlers.Embedder` and `lifecycle.WarmupModel` both existed.
- `defaultState` and `DefaultState()` existed.

## Fix

- Removed inactive ONNX-only fields from `handlers.RuntimeHealth`.
- Removed stale ONNX runtime health test.
- Added shared `embed.Embedder`.
- Updated handlers, batch helpers, lifecycle warmup, warmup handler, main, and tests to use `embed.Embedder`.
- Removed `defaultState` and `DefaultState()` from lifecycle.
- Updated main to use explicit `lifecycle.NewState()`.

## Green Evidence

Commands:

```bash
cd api && go test ./handlers ./lifecycle
cd api && go test ./...
```

Results:

```text
go test ./handlers ./lifecycle: 101 passed in 2 packages
go test ./...: 280 passed in 10 packages
```

Post-cleanup static proof:

```text
gsd_exec d75568af-277e-40e2-a28b-e6ee373d28dd
PASS M048 S02 runtime contract cleanup static proof
```

## Requirement Outcome

- R038 validated for issue #7 findings #26, #29, and #30.

## Residual Issue #7 Findings

Deferred to S03:

- #24 malformed validation message.
- #31 OpenAPI helper silent key drop.
