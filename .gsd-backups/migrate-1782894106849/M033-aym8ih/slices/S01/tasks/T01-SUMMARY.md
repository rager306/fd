---
id: T01
parent: S01
milestone: M033-aym8ih
key_files:
  - tools/provision_onnx_artifacts.py
key_decisions:
  - ONNX Runtime zip/wheel extraction is inferred for `.whl`/`.zip` sources when `source_contract.onnx_runtime.library_member` is present; direct-file runtime sources remain supported for non-zip inputs.
duration: 
verification_result: passed
completed_at: 2026-05-21T07:14:24.534Z
blocker_discovered: false
---

# T01: Implemented safe ONNX Runtime wheel extraction in the provisioning helper.

**Implemented safe ONNX Runtime wheel extraction in the provisioning helper.**

## What Happened

Updated the provisioning helper to read ONNX Runtime member, size, and sha from `source_contract.onnx_runtime`, add safe zip/wheel member extraction, and use the manifest runtime sha when `--onnx-runtime-sha256` is omitted. The helper still copies direct runtime library sources when the source is not `.zip`/`.whl`, preserving direct-file fallback. A synthetic positive probe showed the helper extracts a wheel member to `.gsd/runtime/onnxruntime/libonnxruntime.so.1.26.0` and verifies its checksum.

## Verification

`python3 -m py_compile tools/provision_onnx_artifacts.py` passed; synthetic wheel provisioning probe passed with runtime sha match.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py` | 0 | ✅ pass — no compile errors | 8800ms |
| 2 | `gsd_exec M033 synthetic positive ONNX Runtime wheel provisioning probe` | 0 | ✅ pass — wheel member extracted and sha matched | 147ms |

## Deviations

None.

## Known Issues

Synthetic probe uses a small fake wheel; final guardrails should also run dry-run/verifier and project checks. The exact ONNX model binary remains unhosted.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py`
