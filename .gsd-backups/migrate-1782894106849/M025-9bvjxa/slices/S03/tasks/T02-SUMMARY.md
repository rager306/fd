---
id: T02
parent: S03
milestone: M025-9bvjxa
key_files:
  - docs/onnx-artifacts/README.md
  - docs/onnx-artifacts/PROVISIONING.md
key_decisions:
  - Docs state the manual workflow is `workflow_dispatch` only and artifact sources must be explicit/non-secret or future masked secrets.
  - Docs warn not to pass signed or secret-bearing URLs as plain workflow inputs.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:47:12.061Z
blocker_discovered: false
---

# T02: Documented the manual hosted ONNX CI workflow and URL safety guidance.

**Documented the manual hosted ONNX CI workflow and URL safety guidance.**

## What Happened

Updated the ONNX README and provisioning contract to document the manual hosted CI workflow. The docs explain required inputs, workflow scope, non-secret URL guidance, and why the workflow is not push/PR triggered. Actionlint still passes for both workflows.

## Verification

Docs reference the workflow and safety guidance; actionlint passes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `grep workflow and URL safety docs` | 0 | ✅ pass — manual_workflow_docs=pass | 0ms |
| 2 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/onnx-packaging.yml .github/workflows/go-quality.yml` | 0 | ✅ pass — actionlint no findings | 8400ms |

## Deviations

None.

## Known Issues

Secret-backed artifact URLs are not wired into the workflow yet; current skeleton is safest with non-secret immutable URLs or future hardened secret wiring.

## Files Created/Modified

- `docs/onnx-artifacts/README.md`
- `docs/onnx-artifacts/PROVISIONING.md`
