---
id: T01
parent: S01
milestone: M020-mvkq4d
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
key_decisions:
  - Update the existing `user-bge-m3-dense-fp32.json` manifest rather than creating a separate file, because the ONNX binary is the same dynamic-axis artifact and the new information is validated runtime contract, not a new exported binary.
  - Keep `export.sequence_length=128` as export provenance and add separate validated runtime fields for 1024.
duration: 
verification_result: passed
completed_at: 2026-05-20T09:59:23.870Z
blocker_discovered: false
---

# T01: Chose to update the existing manifest with a separate validated runtime contract.

**Chose to update the existing manifest with a separate validated runtime contract.**

## What Happened

Inspected the current ONNX manifest. It already records dynamic sequence axes and export-time sequence length 128. The correct metadata shape is to preserve export provenance and add a separate validated runtime contract documenting max sequence length 1024, evidence artifacts, and remaining gates.

## Verification

Read the full manifest and attempted GitNexus impact; JSON manifest is not indexed as a symbol, so validation will be JSON/field based.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `read docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | 0 | ✅ pass — manifest inspected | 0ms |
| 2 | `gitnexus_impact target=docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | 0 | ⚠️ not indexed — JSON has no symbol target; use artifact validation | 0ms |

## Deviations

None.

## Known Issues

GitNexus impact cannot target the JSON manifest as a symbol; follow-up verification will use JSON validation and text reference checks.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
