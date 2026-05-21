---
id: T01
parent: S02
milestone: M033-aym8ih
key_files:
  - docs/onnx-artifacts/PROVISIONING.md
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T07:17:16.881Z
blocker_discovered: false
---

# T01: Documented ONNX Runtime wheel extraction support.

**Documented ONNX Runtime wheel extraction support.**

## What Happened

Updated provisioning documentation with the new ONNX Runtime wheel extraction behavior: when a `.whl`/`.zip` source is supplied and the manifest has `source_contract.onnx_runtime.library_member`, the helper extracts only that member, enforces size, writes the runtime library destination, and verifies sha. The docs also preserve direct-file fallback and the remaining hosted/production blockers.

## Verification

Docs contain wheel extraction behavior, source_contract fields, no-extractall boundary, direct-file fallback, and no-production-readiness boundary.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_exec M033 docs wheel extraction checks` | 0 | ✅ pass — required documentation markers present | 64ms |

## Deviations

None.

## Known Issues

Docs still clearly state hosted CI and ONNX production readiness are not implied.

## Files Created/Modified

- `docs/onnx-artifacts/PROVISIONING.md`
