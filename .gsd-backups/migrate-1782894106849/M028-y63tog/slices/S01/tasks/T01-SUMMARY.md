---
id: T01
parent: S01
milestone: M028-y63tog
key_files:
  - api/main.go
  - api/embed/onnx_manifest.go
  - api/handlers/health.go
  - tools/provision_onnx_artifacts.py
  - tools/verify_onnx_artifacts.py
  - .github/workflows/onnx-packaging.yml
  - tools/build_onnx_image.sh
  - docs/onnx-artifacts/OPERATIONS.md
key_decisions:
  - Review scope is read-only and covers manual workflow inputs, provisioning downloads/copies, manifest path resolution, Go startup preflight, health metadata, verifier, Docker packaging, and operations docs.
duration: 
verification_result: passed
completed_at: 2026-05-21T04:14:48.159Z
blocker_discovered: false
---

# T01: Mapped ONNX operational security attack surfaces for review.

**Mapped ONNX operational security attack surfaces for review.**

## What Happened

Mapped ONNX operational security attack surfaces with file:line snippets. Entry points include ONNX env vars in Go startup, manual workflow_dispatch inputs, explicit URL/file sources in the provisioning helper, manifest-local artifact paths, Docker packaging copy paths, verifier path resolution, startup logs/errors, and `/health.runtime` metadata. Trust boundaries are mostly local/admin/operator controlled, with GitHub workflow dispatch and remote artifact URLs as the highest-risk surface.

## Verification

Line-numbered attack surface snippets were extracted and reviewed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_exec M028 attack surface snippet extraction with line numbers` | 0 | ✅ pass — snippets extracted to .gsd/exec/ce7b4851-b4aa-4fbb-8243-d7c2b5d53eca.stdout | 64ms |
| 2 | `read tools/provision_onnx_artifacts.py, .github/workflows/onnx-packaging.yml, tools/build_onnx_image.sh, tools/verify_onnx_artifacts.py, api/main.go snippets` | 0 | ✅ pass — concrete sources/sinks reviewed | 0ms |

## Deviations

None.

## Known Issues

Potential findings are being reported in T02 rather than remediated in this milestone.

## Files Created/Modified

- `api/main.go`
- `api/embed/onnx_manifest.go`
- `api/handlers/health.go`
- `tools/provision_onnx_artifacts.py`
- `tools/verify_onnx_artifacts.py`
- `.github/workflows/onnx-packaging.yml`
- `tools/build_onnx_image.sh`
- `docs/onnx-artifacts/OPERATIONS.md`
