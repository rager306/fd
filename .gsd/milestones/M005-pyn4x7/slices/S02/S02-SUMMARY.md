---
id: S02
parent: M005-pyn4x7
milestone: M005-pyn4x7
provides:
  - Discoverable operational hardening guidance for future deployment/optimization work.
requires:
  []
affects:
  - S03
key_files:
  - README.md
  - .gsd/DECISIONS.md
key_decisions:
  - D001: Current TEI CPU/Candle fallback runtime remains acceptable; ONNX export is future measured optimization requiring A/B evidence.
patterns_established:
  - Runtime warnings should be classified as correctness bugs, host deployment notes, or measured optimization candidates.
observability_surfaces:
  - README Operational Notes
  - DECISIONS.md D001
drill_down_paths:
  - .gsd/milestones/M005-pyn4x7/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M005-pyn4x7/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M005-pyn4x7/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T10:43:31.286Z
blocker_discovered: false
---

# S02: Runtime hardening notes

**S02 made Redis/TEI/logging hardening notes durable in README and DECISIONS.md.**

## What Happened

S02 documented operational hardening notes and recorded the TEI backend decision. README now explains Redis localhost binding, warns against exposing Redis broadly, documents Redis overcommit as host-level tuning, describes TEI ONNX/Candle fallback status, and documents LOG_LEVEL debug cache events. GSD decision D001 records that ONNX export is a future measured optimization rather than an immediate correctness fix.

## Verification

All S02 tasks complete and verified.

## Requirements Advanced

- Operability documentation improved. — 

## Requirements Validated

- Redis binding rationale documented. — 
- Redis overcommit note documented. — 
- TEI ONNX decision recorded. — 
- LOG_LEVEL/debug cache behavior documented. — 

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

README documents host-level Redis overcommit tuning but does not apply host sysctl changes; that remains deployment-specific.

## Follow-ups

S03 should run final verification and commit the documentation/decision work locally.

## Files Created/Modified

- `README.md` — Added operational hardening notes for Redis, TEI, and runtime logging.
- `.gsd/DECISIONS.md` — Recorded TEI ONNX fallback decision.
