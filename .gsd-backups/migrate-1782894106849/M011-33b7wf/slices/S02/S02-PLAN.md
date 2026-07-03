# S02: Opt in backend seam

**Goal:** Introduce the smallest runtime backend selection and manifest validation seam needed for opt-in ONNX, while preserving TEI as the default and avoiding ONNX inference wiring in this slice.
**Demo:** After this, the Go service has a backend selection seam with TEI still default and ONNX rejected or disabled unless explicitly configured.

## Must-Haves

- TEI remains default when no new config is set.
- Backend config validation is explicit and tested.
- ONNX manifest validation checks path, size, SHA256, output name, and dimensions.
- Invalid ONNX config fails fast with actionable error.
- API handlers/cache behavior remain unchanged.
- GitNexus impact is run before code edits.

## Proof Level

- This slice proves: Go tests prove TEI default behavior, backend config validation, and artifact manifest validation failure modes.

## Integration Closure

Provides tested config/manifest validation primitives consumed by S03 ONNX loader. API handlers/cache behavior remain unchanged.

## Verification

- Startup/config errors can identify invalid backend, missing manifest, missing artifact, checksum mismatch, and metadata mismatch without logging raw input text or secrets.

## Tasks

- [x] **T01: Inspect backend seam and impact** `est:small`
  Inspect Go startup/config wiring and cache/embedder seams. Run GitNexus impact analysis on candidate symbols before edits. Decide where backend config and manifest validation belong with minimal package churn.
  - Files: `api/main.go`, `api/embed/`
  - Verify: GitNexus impact recorded for symbols to edit; target files and tests identified.

- [x] **T02: Implement ONNX manifest validation** `est:medium`
  Add a small ONNX artifact manifest type and validation function in the Go API. It should parse the tracked manifest schema subset, resolve local artifact path, check file existence, size, SHA256, output name `dense_vecs`, expected dimensions `1024`, and `production_default=false`. Include unit tests for valid manifest, missing file, checksum mismatch, invalid output, and invalid dimensions.
  - Files: `api/embed/onnx_manifest.go`, `api/embed/onnx_manifest_test.go`
  - Verify: `cd api && go test ./embed -run 'Test.*ONNX.*|Test.*Manifest.*'` passes.

- [x] **T03: Add opt in backend config validation** `est:medium`
  Add runtime backend config parsing in startup code with defaults preserving TEI. Supported backend values should be `tei` and `onnx`; empty/default resolves to `tei`. When `EMBEDDING_BACKEND=onnx`, require `ONNX_ARTIFACT_MANIFEST` and validate it using the manifest validator, but do not wire ONNX inference yet. Tests should prove default TEI, invalid backend, missing manifest, and invalid manifest behavior.
  - Files: `api/main.go`, `api/main_test.go`
  - Verify: `cd api && go test ./... -short` passes; default startup config remains TEI.

- [x] **T04: Verify backend seam safety** `est:small`
  Run S02 verification gates: Go tests, pinned lint, manifest validation against the local M010 artifact, Docker Compose config, and GitNexus detect_changes. Record whether any production behavior changed.
  - Verify: Go tests/lint/config/GitNexus pass; TEI remains default.

## Files Likely Touched

- api/main.go
- api/embed/
- api/embed/onnx_manifest.go
- api/embed/onnx_manifest_test.go
- api/main_test.go
