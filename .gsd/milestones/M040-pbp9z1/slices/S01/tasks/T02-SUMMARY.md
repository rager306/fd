---
id: T02
parent: S01
milestone: M040-pbp9z1
key_files:
  - api/handlers/health.go (RuntimeHealth struct, HealthHandler, NewHealthHandler with TEI-safe metadata)
  - api/handlers/health_test.go (TEI and ONNX health metadata tests)
  - api/main.go (embeddingRuntimeConfig.Health() method wiring TEI and ONNX safe metadata)
  - docs/same-host-embedding-service-contract.md (new canonical contract for same-host clients)
  - README.md (added contract link at line 59)
key_decisions:
  - TEI /health omits runtime block by design — absence of runtime block = TEI backend active; when TEI backend is active and runtime is configured, the block is present with safe fields only
  - ONNX opt-in adds full runtime metadata block with artifact/tokenizer/runtime SHA256 verification fields as pointer bools (omitted from JSON when nil)
  - Model Dimensions fixed at 1024 for deepvk/USER-bge-m3 TEI backend — no inference probe needed, consistent with contract doc
duration: 
verification_result: passed
completed_at: 2026-05-22T04:52:09.466Z
blocker_discovered: false
---

# T02: TEI runtime metadata now exposed via /health with safe fields: backend, model, dimensions, production_default, cache_namespace — no paths, tokens, or secrets

**TEI runtime metadata now exposed via /health with safe fields: backend, model, dimensions, production_default, cache_namespace — no paths, tokens, or secrets**

## What Happened

Updated the runtime health path so the TEI/default backend reports safe runtime metadata comparable to ONNX basics. `api/handlers/health.go` now exposes a `RuntimeHealth` struct with backend, model, dimensions (1024 fixed for deepvk/USER-bge-m3), production_default flag, and cache_namespace for TEI — omitting ONNX-only fields (artifact_id, provider, *_verified pointer fields) to avoid overclaiming. The `Health()` method on `embeddingRuntimeConfig` in `api/main.go` returns this safe TEI block; ONNX path additionally populates the pointer-bool verification fields. `api/handlers/health_test.go` was updated with `TestNewHealthHandlerIncludesSafeTEIRuntimeMetadata` and `TestNewHealthHandlerIncludesSafeRuntimeMetadata` covering field presence/absence, type checks, and secret-leak guards. The same-host embedding service contract doc (`docs/same-host-embedding-service-contract.md`) was created and linked from README.md at line 59.

## Verification

Go tests: `cd api && go test ./... -short` — all packages pass. Linter: `golangci-lint` — 0 issues. Secret-leak check: no signed URLs, tokens, or secret references in health.go, health_test.go, or main.go. Contract doc: 16 KB, all required sections present. README link: confirmed at line 59.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go test ./... -short` | 0 | ✅ pass | 214ms |
| 2 | `cd /root/fd/api && golangci-lint run` | 0 | ✅ pass | 1962ms |
| 3 | `grep -n 'signed|token=|X-Amz|BEGIN.*PRIVATE' health.go health_test.go main.go` | 1 | ✅ pass (no matches = clean) | 13ms |

## Deviations

None

## Known Issues

None

## Files Created/Modified

- `api/handlers/health.go (RuntimeHealth struct, HealthHandler, NewHealthHandler with TEI-safe metadata)`
- `api/handlers/health_test.go (TEI and ONNX health metadata tests)`
- `api/main.go (embeddingRuntimeConfig.Health() method wiring TEI and ONNX safe metadata)`
- `docs/same-host-embedding-service-contract.md (new canonical contract for same-host clients)`
- `README.md (added contract link at line 59)`
