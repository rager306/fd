---
id: T01
parent: S01
milestone: M030-3kdha1
key_files:
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
key_decisions:
  - Go ONNX manifest `artifact.local_path` is now repo-relative and approved-root constrained. Existing `.gsd/runtime/onnx/...` layout remains allowed.
duration: 
verification_result: passed
completed_at: 2026-05-21T05:23:29.050Z
blocker_discovered: false
---

# T01: Hardened Go ONNX manifest artifact path policy and diagnostics.

**Hardened Go ONNX manifest artifact path policy and diagnostics.**

## What Happened

Implemented Go manifest path policy and safer diagnostics. `artifact.local_path` now rejects absolute paths, traversal outside the repository, and unapproved roots. Error messages now use the manifest path value/safe display instead of resolved absolute host paths where possible. Tests now use the existing approved `.gsd/runtime/onnx/...` layout and cover repo-external and unapproved-root rejection. Targeted manifest tests pass.

## Verification

Targeted ONNX manifest tests passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gofmt -w api/embed/onnx_manifest.go api/embed/onnx_manifest_test.go && cd api && go test ./embed -run 'TestValidateONNXArtifactManifest|TestLoadONNXArtifactManifest' -count=1` | 0 | ✅ pass — fd-api/embed ok | 0ms |

## Deviations

None.

## Known Issues

Python verifier/provisioning path policy is still handled in T02.

## Files Created/Modified

- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`
