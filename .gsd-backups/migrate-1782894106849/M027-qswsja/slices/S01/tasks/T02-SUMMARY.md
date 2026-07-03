---
id: T02
parent: S01
milestone: M027-qswsja
key_files:
  - api/main.go
  - api/main_test.go
  - api/handlers/health.go
  - api/handlers/health_test.go
key_decisions:
  - ONNX runtime library sha verification is explicit opt-in through `ONNX_RUNTIME_SHA256` to avoid mandatory startup hashing until runtime artifact source is tracked.
  - Only `CPUExecutionProvider` is accepted by current Go ONNX startup config because current runtime construction uses the CPU execution path.
duration: 
verification_result: passed
completed_at: 2026-05-20T12:38:10.504Z
blocker_discovered: false
---

# T02: Implemented tokenizer, runtime library, and provider startup preflight diagnostics.

**Implemented tokenizer, runtime library, and provider startup preflight diagnostics.**

## What Happened

Implemented startup preflight for tokenizer JSON metadata, optional ONNX Runtime library sha, and provider configuration. `loadEmbeddingRuntimeConfig` now validates tokenizer JSON size/sha when the manifest provides it, checks `ONNX_RUNTIME_SHA256` when supplied, rejects unsupported providers, and populates safe runtime health fields for provider/tokenizer/runtime verification. Added targeted tests for checksum mismatch, unsupported provider, runtime sha success/mismatch, and health metadata.

## Verification

Targeted main/health/manifest tests passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gofmt -w api/embed/onnx_manifest.go api/embed/onnx_manifest_test.go api/main.go api/main_test.go api/handlers/health.go api/handlers/health_test.go && cd api && go test ./embed ./handlers . -run 'TestValidateONNXArtifactManifest|TestLoadEmbeddingRuntimeConfig|TestEmbeddingRuntimeConfigHealth|TestHealth' -count=1` | 0 | ✅ pass — fd-api/embed ok, fd-api/handlers ok, fd-api ok | 0ms |

## Deviations

Provider diagnostics validate configured provider support; they do not enumerate runtime provider availability because the current Go path does not expose that as a startup surface.

## Known Issues

Startup errors may still include some filesystem-path context inherited from existing artifact validation; a dedicated security/logging review remains planned.

## Files Created/Modified

- `api/main.go`
- `api/main_test.go`
- `api/handlers/health.go`
- `api/handlers/health_test.go`
