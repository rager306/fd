# S01: Startup diagnostics and health metadata

**Goal:** Implement runtime diagnostics and health metadata surfaces for ONNX opt-in mode.
**Demo:** After this, ONNX startup failures and health output have safe, actionable diagnostic metadata while TEI health remains unchanged.

## Must-Haves

- Runtime status struct exists and redacts sensitive fields.
- Health handler supports default-compatible and metadata modes.
- Main wires runtime status into `/health`.
- ONNX config validates max sequence length against manifest runtime contract.
- Tests cover default health, ONNX metadata, and config errors.

## Proof Level

- This slice proves: Unit tests and guardrail commands.

## Integration Closure

Provides code-level startup/health implementation of the M025 operations contract.

## Verification

- Safe runtime metadata in health/logs and actionable preflight errors.

## Tasks

- [x] **T01: Implement runtime health metadata** `est:medium`
  Add safe runtime status data model and health handler option so `/health` can include runtime metadata while default-compatible behavior remains unchanged when no metadata is supplied.
  - Files: `api/handlers/health.go`, `api/handlers/health_test.go`
  - Verify: Health handler tests pass and default response still has status/time.

- [x] **T02: Implement ONNX startup preflight diagnostics** `est:medium`
  Extend ONNX manifest/runtime config to include validated max sequence length and safe runtime status; fail when configured sequence length exceeds manifest validated contract; log safe metadata.
  - Files: `api/main.go`, `api/main_test.go`, `api/embed/onnx_manifest.go`, `api/embed/onnx_manifest_test.go`
  - Verify: Config/manifest tests pass, including sequence length mismatch and safe status fields.

- [x] **T03: Verify diagnostics implementation** `est:medium`
  Run S01 verification: targeted tests, default Go tests, tagged tests, lint, default Docker build, binary hygiene, cleanup, GitNexus scope.
  - Verify: All S01 checks pass.

## Files Likely Touched

- api/handlers/health.go
- api/handlers/health_test.go
- api/main.go
- api/main_test.go
- api/embed/onnx_manifest.go
- api/embed/onnx_manifest_test.go
