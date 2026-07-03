---
id: T01
parent: S02
milestone: M034-9tfz77
key_files:
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T07:55:57.555Z
blocker_discovered: false
---

# T01: Documented safe manual workflow dispatch inputs and remaining blockers.

**Documented safe manual workflow dispatch inputs and remaining blockers.**

## What Happened

Updated provisioning docs and the ONNX artifacts README with the manual hosted workflow input contract. The docs list required and optional inputs, state that runtime sha is an optional override when manifest sha exists, prohibit signed/plain secret URLs, require explicit user approval before dispatch, and keep the exact ONNX model binary source as the blocking required input.

## Verification

Docs contain required/optional input policy and no production readiness overclaim.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_exec M034 hosted workflow input docs checks final` | 0 | ✅ pass — required workflow input markers present and production-ready overclaim absent | 42ms |

## Deviations

Initial docs check used a forbidden phrase that matched the existing negated sentence `does not mean ... ONNX is production-ready`; reran with a precise check that requires the negation.

## Known Issues

Docs still block hosted proof until exact ONNX binary source exists and explicit user approval is given for workflow dispatch.

## Files Created/Modified

- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`
