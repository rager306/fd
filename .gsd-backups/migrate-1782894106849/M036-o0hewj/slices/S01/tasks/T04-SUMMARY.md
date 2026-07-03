---
id: T04
parent: S01
milestone: M036-o0hewj
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T09:37:11.148Z
blocker_discovered: false
---

# T04: Verified the planned reproducible-export contract with docs/manifest guardrails.

**Verified the planned reproducible-export contract with docs/manifest guardrails.**

## What Happened

Ran S01 verification for the reproducible-export contract. Manifest JSON is valid, tooling compiles, provisioning dry-run and verifiers pass, actionlint passes, custom contract/leak checks pass, and GitNexus reports low-risk docs/manifest changes with no affected processes.

## Verification

All S01 checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `json.tool + py_compile + provisioning dry-run + verify_onnx_artifacts + verify_onnx_export_contract` | 0 | ✅ pass — manifest valid and provisioning/export checks passed | 7100ms |
| 2 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 7000ms |
| 3 | `custom M036 slice contract checks` | 0 | ✅ pass — failed_checks=0, leak_markers=0, signed_url_like=0 | 7000ms |
| 4 | `gitnexus_detect_changes(scope=all)` | 0 | ✅ pass — low risk, no affected processes | 0ms |
| 5 | `bg_shell list` | 0 | ✅ pass — no background processes | 0ms |

## Deviations

None.

## Known Issues

No export regenerated; no source blocker resolved.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`
- `benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt`
