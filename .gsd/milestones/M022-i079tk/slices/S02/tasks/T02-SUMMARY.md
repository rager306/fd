---
id: T02
parent: S02
milestone: M022-i079tk
key_files:
  - docs/onnx-artifacts/README.md
  - .gsd/DECISIONS.md
key_decisions:
  - D020 recorded: default CI gets artifact-free ONNX contract checks; full ONNX image CI is deferred until artifact provisioning/cache exists.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:45:47.295Z
blocker_discovered: false
---

# T02: Documented and recorded the ONNX image CI provisioning blocker.

**Documented and recorded the ONNX image CI provisioning blocker.**

## What Happened

Documented the CI boundary in the ONNX artifact README and recorded D020. The regular workflow can validate manifest metadata and binary hygiene today, but full ONNX image CI remains blocked on provisioning the ONNX model, native tokenizer library, tokenizer JSON, and ONNX Runtime shared library.

## Verification

D020 was saved and README now distinguishes CI-safe contract checks from future full ONNX image CI requirements.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_decision_save` | 0 | ✅ pass — Saved decision D020 | 0ms |
| 2 | `read docs/onnx-artifacts/README.md` | 0 | ✅ pass — CI boundary section present | 0ms |

## Deviations

None.

## Known Issues

Full ONNX image CI still needs an external artifact source/cache before it can be enabled truthfully.

## Files Created/Modified

- `docs/onnx-artifacts/README.md`
- `.gsd/DECISIONS.md`
