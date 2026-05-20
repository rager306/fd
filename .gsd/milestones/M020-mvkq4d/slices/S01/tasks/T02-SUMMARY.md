---
id: T02
parent: S01
milestone: M020-mvkq4d
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
key_decisions:
  - Manifest status is now `experimental_validated_quality_performance_local`, still with `production_default=false`.
  - `export.sequence_length` remains 128 as provenance; `runtime.validated_max_sequence_length` records 1024 as validated runtime contract.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:00:22.804Z
blocker_discovered: false
---

# T02: Updated the ONNX manifest with the validated 1024 runtime contract.

**Updated the ONNX manifest with the validated 1024 runtime contract.**

## What Happened

Updated the tracked ONNX manifest with an explicit 1024 runtime contract. The manifest now records validated runtime env, M018 legal quality PASS metrics, M019 local performance metrics, failure contract entries for sequence length and production misuse, and future production gates. It keeps export sequence length 128 and production_default false.

## Verification

JSON parsing and required field assertions passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m json.tool docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | 0 | ✅ pass — valid JSON | 0ms |
| 2 | `python manifest contract assertions` | 0 | ✅ pass — onnx_manifest_1024_contract=pass | 0ms |

## Deviations

None.

## Known Issues

Manifest remains a local/prototype artifact contract; packaging/CI and external artifact distribution are still future gates.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
