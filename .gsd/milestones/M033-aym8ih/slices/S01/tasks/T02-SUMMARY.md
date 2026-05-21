---
id: T02
parent: S01
milestone: M033-aym8ih
key_files:
  - tools/provision_onnx_artifacts.py
  - .gsd/exec/95eaefb4-9f1f-494d-b22a-d179785ea722.stdout
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T07:15:10.130Z
blocker_discovered: false
---

# T02: Verified ONNX Runtime wheel extraction failure modes and direct-file fallback.

**Verified ONNX Runtime wheel extraction failure modes and direct-file fallback.**

## What Happened

Ran synthetic probes in temporary repo roots. Missing wheel member, oversized member, and checksum mismatch all failed as expected with explicit errors. A direct runtime shared-library source still provisions successfully, preserving fallback behavior for non-wheel sources.

## Verification

Synthetic probes passed: `missing_member`, `oversized_member`, `checksum_mismatch`, and `direct_file_fallback`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_exec M033 synthetic negative and direct fallback probes` | 0 | ✅ pass — all expected failure/success probes behaved correctly | 449ms |

## Deviations

None.

## Known Issues

Synthetic probes cover zip behavior and direct fallback. Full hosted workflow proof remains blocked by exact ONNX model binary source and no push/dispatch.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py`
- `.gsd/exec/95eaefb4-9f1f-494d-b22a-d179785ea722.stdout`
