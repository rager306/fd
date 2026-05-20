---
id: T02
parent: S01
milestone: M026-ji0i9y
key_files:
  - api/main.go
  - api/main_test.go
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
  - api/handlers/health.go
  - api/handlers/health_test.go
key_decisions:
  - Manifest runtime contract now includes `validated_max_sequence_length`; configs fail when `ONNX_MAX_SEQUENCE_LENGTH` exceeds it.
  - Main logs safe ONNX preflight metadata and wires ONNX runtime health metadata only when ONNX is active.
  - TEI/default health route remains default-compatible because `runtimeConfig.Health(...)` returns nil for TEI.
duration: 
verification_result: passed
completed_at: 2026-05-20T12:04:31.986Z
blocker_discovered: false
---

# T02: Implemented ONNX startup preflight diagnostics and runtime health wiring.

**Implemented ONNX startup preflight diagnostics and runtime health wiring.**

## What Happened

Implemented ONNX startup preflight diagnostics. The manifest parser now carries `validated_max_sequence_length`, runtime config rejects sequence lengths above that contract, main logs safe ONNX preflight metadata, Redis connection logs the cache namespace, and `/health` includes safe runtime metadata only when ONNX is active. Targeted tests cover health shape, manifest validation, sequence length mismatch, and runtime health metadata.

## Verification

Targeted handler/embed/main tests passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gofmt -w api/handlers/health.go api/handlers/health_test.go api/main.go api/main_test.go api/embed/onnx_manifest.go api/embed/onnx_manifest_test.go && cd api && go test ./handlers ./embed -run 'TestHealth|TestValidateONNXArtifactManifest' -count=1 && go test . -run 'TestLoadEmbeddingRuntimeConfig|TestEmbeddingRuntimeConfigHealth' -count=1` | 0 | ✅ pass — targeted diagnostics tests passed | 0ms |

## Deviations

None.

## Known Issues

Health metadata exposes safe runtime fields only; deeper code-level startup preflight for tokenizer JSON checksum and ONNX Runtime library hash remains future work.

## Files Created/Modified

- `api/main.go`
- `api/main_test.go`
- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`
- `api/handlers/health.go`
- `api/handlers/health_test.go`
