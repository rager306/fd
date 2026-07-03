# S01: Startup artifact and provider preflight

**Goal:** Implement and test remaining ONNX startup preflight diagnostics.
**Demo:** After this, ONNX startup has stricter artifact preflight for tokenizer JSON, optional runtime library sha, and provider config.

## Must-Haves

- Manifest exposes tokenizer.json expected size/sha metadata.
- Config validates tokenizer JSON path against manifest metadata.
- Config validates ONNX runtime sha when `ONNX_RUNTIME_SHA256` is set.
- Config rejects unsupported `ONNX_PROVIDER` values.
- Runtime health reports provider and runtime library verification flag safely.

## Proof Level

- This slice proves: Targeted tests and guardrails.

## Integration Closure

Completes tokenizer/runtime/provider preflight for opt-in ONNX startup.

## Verification

- Actionable startup errors and safe health metadata for provider/runtime verification.

## Tasks

- [x] **T01: Expose tokenizer metadata from ONNX manifest** `est:small`
  Extend manifest validation types to expose tokenizer JSON source file size/sha and add tests for metadata parsing.
  - Files: `api/embed/onnx_manifest.go`, `api/embed/onnx_manifest_test.go`
  - Verify: Manifest tests pass.

- [x] **T02: Implement tokenizer runtime provider preflight** `est:medium`
  Implement startup preflight for tokenizer JSON checksum, optional ONNX Runtime sha, and provider validation; extend health metadata and main tests.
  - Files: `api/main.go`, `api/main_test.go`, `api/handlers/health.go`, `api/handlers/health_test.go`
  - Verify: Targeted main/health tests pass.

- [x] **T03: Verify startup preflight diagnostics** `est:medium`
  Run S01 verification: targeted tests, default Go tests, lint, tagged checks, default Docker, binary hygiene, cleanup, GitNexus scope.
  - Verify: All S01 checks pass.

## Files Likely Touched

- api/embed/onnx_manifest.go
- api/embed/onnx_manifest_test.go
- api/main.go
- api/main_test.go
- api/handlers/health.go
- api/handlers/health_test.go
