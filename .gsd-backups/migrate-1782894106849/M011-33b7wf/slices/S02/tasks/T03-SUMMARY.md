---
id: T03
parent: S02
milestone: M011-33b7wf
key_files:
  - api/main.go
  - api/main_test.go
key_decisions:
  - Default `EMBEDDING_BACKEND` is `tei`; empty env preserves existing runtime behavior.
  - Supported backend config values are `tei` and `onnx`; invalid values fail validation.
  - `ONNX_ARTIFACT_MANIFEST` is required only when `EMBEDDING_BACKEND=onnx`.
  - Valid ONNX config is validated but not served in S02; explicit ONNX request fails fast until S03 implements inference wiring.
duration: 
verification_result: passed
completed_at: 2026-05-19T19:04:50.747Z
blocker_discovered: false
---

# T03: Added explicit backend config validation with TEI default preserved and ONNX gated behind manifest validation.

**Added explicit backend config validation with TEI default preserved and ONNX gated behind manifest validation.**

## What Happened

Added runtime backend config parsing to `api/main.go`. `loadEmbeddingRuntimeConfig` defaults to TEI, rejects unsupported backend values, requires `ONNX_ARTIFACT_MANIFEST` when ONNX is requested, and validates the manifest via `embed.ValidateONNXArtifactManifest`. `main` logs the selected backend. If ONNX is requested in S02, it validates the artifact and then fails fast with `onnx backend requested but inference wiring is not implemented yet`, avoiding a silent TEI fallback that would corrupt benchmark evidence. Added tests in `api/main_test.go` covering default TEI, invalid backend, missing ONNX manifest, valid ONNX manifest validation, and invalid ONNX manifest rejection.

## Verification

Fresh verification passed: `gofmt` completed and `cd api && go test ./... -short` passed for all packages.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gofmt -w api/main.go api/main_test.go && cd api && go test ./... -short` | 0 | ✅ pass — api, cache, embed, handlers all ok | 0ms |

## Deviations

S02 `main` branch validates ONNX manifests but exits with a not-implemented error if `EMBEDDING_BACKEND=onnx` is requested, rather than silently falling back to TEI. S03 will replace that branch with the actual loader.

## Known Issues

`EMBEDDING_BACKEND=onnx` is not usable for inference yet by design. It validates the manifest and exits with an implementation-not-ready error in `main`; S03 must wire the ONNX backend or adjust this branch.

## Files Created/Modified

- `api/main.go`
- `api/main_test.go`
