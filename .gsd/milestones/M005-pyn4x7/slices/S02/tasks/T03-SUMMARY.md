---
id: T03
parent: S02
milestone: M005-pyn4x7
key_files:
  - .gsd/DECISIONS.md
  - README.md
key_decisions:
  - D001: Keep current measured TEI CPU/Candle fallback runtime for now; treat ONNX export as future measured optimization requiring A/B benchmark evidence.
duration: 
verification_result: passed
completed_at: 2026-05-19T10:43:05.915Z
blocker_discovered: false
---

# T03: Recorded D001: ONNX export is a measured future optimization, not an immediate runtime fix.

**Recorded D001: ONNX export is a measured future optimization, not an immediate runtime fix.**

## What Happened

Recorded a GSD decision for TEI backend handling. The project will keep the current measured TEI CPU/Candle fallback runtime for now and treat ONNX export/runtime changes as a future measured optimization requiring A/B benchmark evidence before becoming the default. This prevents future agents from interpreting ONNX-missing warnings as an immediate correctness bug.

## Verification

`gsd_decision_save` succeeded and regenerated `.gsd/DECISIONS.md`.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_decision_save(scope=runtime-performance, decision=TEI ONNX fallback handling)` | 0 | ✅ pass: saved decision D001 | 0ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `.gsd/DECISIONS.md`
- `README.md`
