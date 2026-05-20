---
id: T01
parent: S03
milestone: M025-9bvjxa
key_files:
  - .github/workflows/onnx-packaging.yml
key_decisions:
  - Full ONNX hosted CI is a manual `workflow_dispatch` workflow, not push/PR required CI.
  - Workflow requires explicit ONNX and native tokenizer sources and fails early if missing.
  - Workflow separates provisioning, verification, tagged tests, image build, and next-gate summary.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:46:07.632Z
blocker_discovered: false
---

# T01: Added a safe manual ONNX packaging workflow skeleton.

**Added a safe manual ONNX packaging workflow skeleton.**

## What Happened

Added `.github/workflows/onnx-packaging.yml` as a manual workflow skeleton. It accepts explicit artifact inputs, validates required inputs, provisions artifacts via the new helper, verifies them strictly, runs tagged tokenizer and ONNX smoke tests, builds the opt-in ONNX image, and reminds operators that legal/performance gates are still required before production promotion. `actionlint` passed.

## Verification

Workflow path was new and actionlint passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test ! -e .github/workflows/onnx-packaging.yml` | 0 | ✅ pass — onnx_workflow_path_new=pass | 0ms |
| 2 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/onnx-packaging.yml` | 0 | ✅ pass — actionlint no findings | 10500ms |

## Deviations

None.

## Known Issues

Workflow cannot be run truthfully until immutable artifact URLs or equivalent accessible sources are supplied.

## Files Created/Modified

- `.github/workflows/onnx-packaging.yml`
