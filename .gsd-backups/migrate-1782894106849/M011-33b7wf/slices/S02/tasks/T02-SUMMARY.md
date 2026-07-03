---
id: T02
parent: S02
milestone: M011-33b7wf
key_files:
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
key_decisions:
  - Manifest validation lives in `api/embed` as a pure Go dependency-free validator.
  - Validation checks production_default, git_tracked, artifact path, size, SHA256, output name, dimensions, and normalization expectation before future ONNX load.
  - Validation returns sentinel errors for missing artifact, checksum mismatch, metadata mismatch, and production-default violations.
duration: 
verification_result: mixed
completed_at: 2026-05-19T19:02:45.861Z
blocker_discovered: false
---

# T02: Implemented and tested the ONNX artifact manifest validator in Go.

**Implemented and tested the ONNX artifact manifest validator in Go.**

## What Happened

Implemented `api/embed/onnx_manifest.go` and tests. The validator parses the tracked manifest schema subset, rejects production-default manifests and git-tracked artifacts, validates local path, size, SHA256, `dense_vecs` output, 1024 dimensions, and normalized-output expectation. Tests cover valid manifest, missing artifact, checksum mismatch, invalid output name, invalid dimensions, production-default rejection, and invalid JSON. Focused embed tests passed after formatting.

## Verification

Fresh verification passed: `gofmt` completed, `cd api && go test ./embed -run 'Test.*ONNX.*|Test.*Manifest.*'` passed with 7 tests. LSP diagnostics were attempted but unavailable for Go.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `lsp diagnostics api/embed/onnx_manifest.go` | 1 | ⚠️ unavailable — No language server found | 0ms |
| 2 | `gofmt -w api/embed/onnx_manifest.go api/embed/onnx_manifest_test.go && cd api && go test ./embed -run 'Test.*ONNX.*|Test.*Manifest.*'` | 0 | ✅ pass — ok fd-api/embed | 0ms |

## Deviations

LSP diagnostics were unavailable for Go (`No language server found`), so Go native `gofmt` and `go test` were used.

## Known Issues

The validator currently reads the whole ONNX artifact for SHA256, which is acceptable for startup/prototype validation but may need caching or manifest stamp logic if startup latency becomes an issue.

## Files Created/Modified

- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`
