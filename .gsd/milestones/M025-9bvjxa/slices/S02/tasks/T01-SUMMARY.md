---
id: T01
parent: S02
milestone: M025-9bvjxa
key_files:
  - docs/onnx-artifacts/OPERATIONS.md
  - docs/onnx-artifacts/README.md
key_decisions:
  - Operational rollout is staged and opt-in; rollback is a backend/image switch back to TEI/default plus health/smoke verification.
  - API health details for ONNX are a future implementation gap, not claimed as complete now.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:42:57.247Z
blocker_discovered: false
---

# T01: Documented ONNX operational diagnostics, rollout, and rollback contract.

**Documented ONNX operational diagnostics, rollout, and rollback contract.**

## What Happened

Wrote the ONNX operational diagnostics and rollout contract and linked it from the README. The runbook defines runtime mode boundaries, startup preflight requirements, actionable failure messages, safe logging and health expectations, rollout stages, rollback steps, and current open gaps.

## Verification

New operations path was unused before writing and README now links the operations contract.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test ! -e docs/onnx-artifacts/OPERATIONS.md` | 0 | ✅ pass — operations_path_new=pass | 0ms |
| 2 | `write docs/onnx-artifacts/OPERATIONS.md; edit docs/onnx-artifacts/README.md` | 0 | ✅ pass — operations contract written and linked | 0ms |

## Deviations

None.

## Known Issues

Operational health details and startup preflight status surfaces are documented but not fully implemented in API responses yet.

## Files Created/Modified

- `docs/onnx-artifacts/OPERATIONS.md`
- `docs/onnx-artifacts/README.md`
