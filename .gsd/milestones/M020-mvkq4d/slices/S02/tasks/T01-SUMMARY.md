---
id: T01
parent: S02
milestone: M020-mvkq4d
key_files:
  - .gsd/DECISIONS.md
  - .gsd/gsd.db
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
key_decisions:
  - D018 recorded: existing manifest is the experimental ONNX 1024 runtime contract; production remains blocked on Docker/CI packaging and artifact provisioning.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:01:41.300Z
blocker_discovered: false
---

# T01: Recorded the ONNX 1024 artifact contract decision.

**Recorded the ONNX 1024 artifact contract decision.**

## What Happened

Recorded D018 in the GSD decision register. The decision captures why the existing dynamic-axis ONNX manifest was updated rather than creating a new binary artifact contract, and preserves the distinction between validated experimental runtime and production readiness.

## Verification

`gsd_decision_save` returned `Saved decision D018`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_decision_save` | 0 | ✅ pass — Saved decision D018 | 0ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `.gsd/DECISIONS.md`
- `.gsd/gsd.db`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
