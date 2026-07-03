---
id: T02
parent: S01
milestone: M036-o0hewj
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T09:34:37.774Z
blocker_discovered: false
---

# T02: Persisted the planned reproducible-export workflow contract without claiming proof.

**Persisted the planned reproducible-export workflow contract without claiming proof.**

## What Happened

Added a planned reproducible-export workflow contract to the ONNX manifest and provisioning docs. The contract records pinned model/toolchain/export inputs, expected artifact contract, acceptance gates, success/failure interpretation, and forbidden claims. It explicitly keeps the current verifier boundary as `existing_artifact_contract_verification_not_regenerated_export` and marks the workflow as `planned_not_proven`.

## Verification

Manifest JSON and reproducible-export contract marker checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m json.tool docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | 0 | ✅ pass — manifest JSON valid | 0ms |
| 2 | `gsd_exec M036 manifest docs reproducible export checks` | 0 | ✅ pass — planned_not_proven, pinned inputs, gates, and forbidden claims present | 93ms |

## Deviations

None.

## Known Issues

The reproducible-export workflow contract is planned only. No ONNX export was regenerated and no proof exists yet.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
