---
id: T02
parent: S02
milestone: M042-fjf2en
key_files:
  - api/main.go
  - api/main_test.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T10:37:01.028Z
blocker_discovered: false
---

# T02: Removed ONNX as an accepted active runtime backend from fd startup config; startup now supports TEI only.

**Removed ONNX as an accepted active runtime backend from fd startup config; startup now supports TEI only.**

## What Happened

Updated `api/main.go` so `embeddingRuntimeConfig` contains only the active TEI backend, `/health.runtime` metadata is TEI-only, and `loadEmbeddingRuntimeConfig` rejects any `EMBEDDING_BACKEND` value other than `tei` with a clear TEI-only error. Removed ONNX manifest/runtime/tokenizer verification from `main.go` and removed ONNX client initialization from startup; fd always constructs the TEI client. Updated `api/main_test.go` to assert stale `ONNX_*` env is ignored for default TEI and `EMBEDDING_BACKEND=onnx` fails closed instead of being validated.

## Verification

Targeted startup/config tests passed: `cd api && go test . -run 'TestLoadEmbeddingRuntimeConfig|TestEmbeddingRuntimeConfigHealth'`. Full module tests also passed after the edit: `cd api && go test ./...`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test . -run 'TestLoadEmbeddingRuntimeConfig|TestEmbeddingRuntimeConfigHealth'` | 0 | ✅ pass: TEI-only runtime config tests pass | 120000ms |
| 2 | `cd api && go test ./...` | 0 | ✅ pass: full module tests pass after TEI-only startup change | 180000ms |

## Deviations

This task removes ONNX from active startup/config only. Build-tagged ONNX package files, docs, CI/tooling references are intentionally left for T03/T04.

## Known Issues

TEI internal ONNX/ORT probing is inside the external TEI container and is not affected by this fd startup cleanup; D048 tracks that as separate TEI startup stabilization.

## Files Created/Modified

- `api/main.go`
- `api/main_test.go`
