---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Add opt in backend config validation

Add runtime backend config parsing in startup code with defaults preserving TEI. Supported backend values should be `tei` and `onnx`; empty/default resolves to `tei`. When `EMBEDDING_BACKEND=onnx`, require `ONNX_ARTIFACT_MANIFEST` and validate it using the manifest validator, but do not wire ONNX inference yet. Tests should prove default TEI, invalid backend, missing manifest, and invalid manifest behavior.

## Inputs

- `api/embed/onnx_manifest.go`

## Expected Output

- `api/main.go`
- `api/main_test.go`

## Verification

`cd api && go test ./... -short` passes; default startup config remains TEI.

## Observability Impact

Startup errors become explicit and test-covered before runtime inference changes.
