---
id: T01
parent: S01
milestone: M027-qswsja
key_files:
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
key_decisions:
  - Expose tokenizer.json source-file metadata from the ONNX manifest validation result without making it a standalone artifact validator.
duration: 
verification_result: passed
completed_at: 2026-05-20T12:37:40.293Z
blocker_discovered: false
---

# T01: Exposed tokenizer JSON size/sha metadata from ONNX manifest validation.

**Exposed tokenizer JSON size/sha metadata from ONNX manifest validation.**

## What Happened

Extended the ONNX manifest model to parse `model.source_files.tokenizer.json` metadata and expose tokenizer JSON size and sha256 through `ONNXArtifactValidation`. Added tests for valid metadata and invalid tokenizer sha metadata. Verified targeted manifest tests pass.

## Verification

Targeted manifest tests passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./embed ./handlers . -run 'TestValidateONNXArtifactManifest|TestLoadEmbeddingRuntimeConfig|TestEmbeddingRuntimeConfigHealth|TestHealth' -count=1` | 0 | ✅ pass — fd-api/embed ok, fd-api/handlers ok, fd-api ok | 0ms |
| 2 | `python3 manifest metadata check` | 0 | ✅ pass — tokenizer_manifest_size=True tokenizer_manifest_sha=True | 0ms |

## Deviations

None.

## Known Issues

Tokenizer file existence/checksum validation is performed by startup config in T02, not by manifest artifact validation alone.

## Files Created/Modified

- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`
