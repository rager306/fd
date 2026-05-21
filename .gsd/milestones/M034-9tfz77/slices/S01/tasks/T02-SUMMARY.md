---
id: T02
parent: S01
milestone: M034-9tfz77
key_files:
  - .github/workflows/onnx-packaging.yml
  - tools/provision_onnx_artifacts.py
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T07:53:01.537Z
blocker_discovered: false
---

# T02: Verified manual workflow and provisioning alignment.

**Verified manual workflow and provisioning alignment.**

## What Happened

Verified local workflow alignment after the YAML change. Actionlint passes. Provisioning dry-run still succeeds and exposes the manifest-derived ONNX Runtime sha, confirming the workflow can omit `onnx_runtime_sha256` while the helper still has a checksum source.

## Verification

actionlint, py_compile, provisioning dry-run, and alignment check passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py && python3 tools/provision_onnx_artifacts.py --dry-run ...` | 0 | ✅ pass — dry-run includes manifest-derived runtime sha | 14400ms |
| 2 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 14400ms |
| 3 | `gsd_exec M034 workflow/provisioning alignment check` | 0 | ✅ pass — runtime sha found in dry-run and workflow marker present | 139ms |

## Deviations

None.

## Known Issues

The workflow still requires ONNX and native tokenizer sources; exact ONNX model source remains blocked. Workflow not dispatched.

## Files Created/Modified

- `.github/workflows/onnx-packaging.yml`
- `tools/provision_onnx_artifacts.py`
