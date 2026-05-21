---
id: T01
parent: S01
milestone: M034-9tfz77
key_files:
  - .github/workflows/onnx-packaging.yml
key_decisions:
  - The workflow now lets provisioning use manifest-derived ONNX Runtime sha when `onnx_runtime_sha256` is omitted, matching M033 helper behavior.
duration: 
verification_result: passed
completed_at: 2026-05-21T07:52:18.870Z
blocker_discovered: false
---

# T01: Aligned manual workflow runtime sha handling with M033 provisioning behavior.

**Aligned manual workflow runtime sha handling with M033 provisioning behavior.**

## What Happened

Updated the manual ONNX packaging workflow so `onnx_runtime_sha256` is an optional override instead of a hard requirement. If a runtime source is provided without the sha input, the validation step logs that provisioning will use `source_contract.onnx_runtime.library_sha256`. The provisioning step now passes `--onnx-runtime-sha256` only when non-empty. Actionlint and workflow text checks passed.

## Verification

actionlint passed and workflow text checks confirmed the manifest-derived sha path.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 6700ms |
| 2 | `gsd_exec M034 workflow runtime sha text checks` | 0 | ✅ pass — required workflow markers present, old hard-fail message absent | 47ms |

## Deviations

None.

## Known Issues

Workflow was not dispatched. Exact ONNX model source remains required before hosted proof.

## Files Created/Modified

- `.github/workflows/onnx-packaging.yml`
