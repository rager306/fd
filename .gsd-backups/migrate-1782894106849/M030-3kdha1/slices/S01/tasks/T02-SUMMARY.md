---
id: T02
parent: S01
milestone: M030-3kdha1
key_files:
  - tools/provision_onnx_artifacts.py
  - tools/verify_onnx_artifacts.py
  - tools/build_onnx_image.sh
key_decisions:
  - Python provisioning/verifier artifact paths are now constrained to approved repo-relative roots.
  - Default tool outputs use repo-relative or basename-safe display; build script missing diagnostics name env keys instead of absolute values.
duration: 
verification_result: passed
completed_at: 2026-05-21T05:27:42.133Z
blocker_discovered: false
---

# T02: Hardened Python provisioning/verifier path policy and diagnostics.

**Hardened Python provisioning/verifier path policy and diagnostics.**

## What Happened

Implemented Python path-root policy and safer path display. Provisioning and verifier now reject absolute, traversal, and unapproved-root manifest artifact paths while preserving approved roots such as `.gsd/runtime/onnx`, `.gsd/runtime/tokenizers`, `.gsd/runtime/onnxruntime`, and `tei-models`. Dry-run/verifier output uses `repo_root: "."` and safe path displays; build script missing diagnostics name the env keys rather than printing full configured paths. Deterministic local probes cover rejected/allowed paths and safe absolute-path display.

## Verification

Python compile, Go manifest tests, and local path policy probes passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py && cd api && go test ./embed -run 'TestValidateONNXArtifactManifest|TestLoadONNXArtifactManifest' -count=1` | 0 | ✅ pass — py_compile and fd-api/embed targeted tests passed | 0ms |
| 2 | `gsd_exec M030 path policy local probes` | 0 | ✅ pass — reject absolute/traversal/unapproved roots; allow approved roots; safe absolute display | 109ms |

## Deviations

None.

## Known Issues

Full provisioning with real external artifact sources remains future work; this only changes local/tool policy and diagnostics.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py`
- `tools/verify_onnx_artifacts.py`
- `tools/build_onnx_image.sh`
