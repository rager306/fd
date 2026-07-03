---
id: T02
parent: S02
milestone: M021-4t2wpt
key_files:
  - .gsd/DECISIONS.md
  - .gsd/gsd.db
key_decisions:
  - D019 recorded: default Docker/CI remains CGO-disabled TEI path; opt-in ONNX runtime requires explicit `onnx hf_tokenizers` tags and verified artifacts.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:23:07.138Z
blocker_discovered: false
---

# T02: Recorded the ONNX Docker/CI packaging boundary decision.

**Recorded the ONNX Docker/CI packaging boundary decision.**

## What Happened

Recorded D019 in the GSD decision register. The decision captures the packaging boundary: default Docker/CI remains independent from ONNX/native artifacts, while ONNX 1024 runtime requires explicit build tags and artifact verification.

## Verification

`gsd_decision_save` returned `Saved decision D019`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_decision_save` | 0 | ✅ pass — Saved decision D019 | 0ms |

## Deviations

None.

## Known Issues

A future dedicated ONNX image target can revisit this once artifact provisioning is implemented.

## Files Created/Modified

- `.gsd/DECISIONS.md`
- `.gsd/gsd.db`
