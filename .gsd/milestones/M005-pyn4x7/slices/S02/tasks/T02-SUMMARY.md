---
id: T02
parent: S02
milestone: M005-pyn4x7
key_files:
  - README.md
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-19T10:42:47.839Z
blocker_discovered: false
---

# T02: Documented Redis exposure, Redis overcommit, TEI ONNX/Candle status, and runtime logging behavior in README.

**Documented Redis exposure, Redis overcommit, TEI ONNX/Candle status, and runtime logging behavior in README.**

## What Happened

Updated README with operational hardening notes. The new section documents why Redis is host-bound only to 127.0.0.1 in the local override, warns against exposing Redis on all interfaces without network controls, explains the Redis `vm.overcommit_memory` host-level warning, clarifies that TEI ONNX export is a future measured optimization rather than a current correctness fix, and documents `LOG_LEVEL` plus debug cache events that avoid raw input text.

## Verification

README now contains Redis localhost binding rationale, overcommit note, TEI ONNX measured-follow-up guidance, and LOG_LEVEL/debug cache notes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `edit README.md configuration and operational notes` | 0 | ✅ pass: hardening notes added | 0ms |

## Deviations

Also documented `LOG_LEVEL` in the configuration table because it was added in M004 but missing from README.

## Known Issues

None.

## Files Created/Modified

- `README.md`
