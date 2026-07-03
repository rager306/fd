---
id: T02
parent: S01
milestone: M037-d23oz4
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T10:14:27.777Z
blocker_discovered: false
---

# T02: Persisted the target-runtime validation contract for Go and future Rust ONNX paths.

**Persisted the target-runtime validation contract for Go and future Rust ONNX paths.**

## What Happened

Added a target-runtime validation contract to the ONNX manifest and provisioning docs. The contract states that Python helper checks are setup/provenance evidence only, defines Go API/package gates for new or regenerated ONNX artifacts, requires Redis namespace isolation, and requires any future Rust runtime to pass its own equivalent gates rather than inheriting Go evidence.

## Verification

Manifest JSON and target-runtime marker checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m json.tool docs/onnx-artifacts/user-bge-m3-dense-fp32.json` | 0 | ✅ pass — manifest JSON valid | 0ms |
| 2 | `gsd_exec M037 manifest docs target runtime checks` | 0 | ✅ pass — Python boundary, Go gates, Rust rule, Redis isolation, and blockers present | 46ms |

## Deviations

None.

## Known Issues

This records the target-runtime validation requirement but does not run new Go/Rust target-runtime gates.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
