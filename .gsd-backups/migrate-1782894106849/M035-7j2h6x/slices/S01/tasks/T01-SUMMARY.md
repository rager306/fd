---
id: T01
parent: S01
milestone: M035-7j2h6x
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T09:19:58.082Z
blocker_discovered: false
---

# T01: Persisted the exact ONNX binary hosting contract without marking a source as available.

**Persisted the exact ONNX binary hosting contract without marking a source as available.**

## What Happened

Added a planned, non-source hosting contract to the ONNX manifest and provisioning docs. The contract records the exact binary size/sha, recommended immutable object key and release filename, allowed and forbidden future source forms, and the pre-dispatch checklist. The manifest still keeps `source_status=blocked` and does not add a `source_url`.

## Verification

JSON validity and contract marker checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m json.tool docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | 0 | ✅ pass — manifest JSON valid | 0ms |
| 2 | `gsd_exec M035 exact binary contract manifest/docs check adjusted` | 0 | ✅ pass — source remains blocked, hosting_contract planned_not_uploaded, no source_url, docs markers present | 47ms |

## Deviations

Initial text check searched the unformatted phrase `not an uploaded source`; docs used markdown emphasis `**not** an uploaded source`. Rechecked with the exact rendered marker.

## Known Issues

No ONNX binary has been uploaded or mirrored; `onnx_source_url` remains blocked.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
