---
id: T03
parent: S01
milestone: M035-7j2h6x
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T09:22:21.287Z
blocker_discovered: false
---

# T03: Verified the exact binary hosting contract and confirmed no external action was taken.

**Verified the exact binary hosting contract and confirmed no external action was taken.**

## What Happened

Ran S01 verification. Manifest JSON is valid, provisioning dry-run still correctly reports missing source blockers, strict artifact/export verifiers pass against local ignored artifacts, actionlint passes, custom contract checks pass with no leak/signed URL markers, and GitNexus reports low-risk docs/manifest changes with no affected processes.

## Verification

All S01 checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m json.tool docs/onnx-artifacts/user-bge-m3-dense-fp32.json && python3 -m py_compile ... && provisioning dry-run && verify_onnx_artifacts && verify_onnx_export_contract` | 0 | ✅ pass — manifest valid, dry-run/verifiers passed | 18200ms |
| 2 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 18100ms |
| 3 | `custom M035 slice contract checks` | 0 | ✅ pass — failed_checks=0, leak_markers=0, signed_url_like=0 | 18100ms |
| 4 | `gitnexus_detect_changes(scope=all)` | 0 | ✅ pass — low risk, no affected processes | 0ms |
| 5 | `bg_shell list` | 0 | ✅ pass — no background processes | 0ms |

## Deviations

None.

## Known Issues

Exact ONNX source remains blocked; no artifact upload, push, or workflow dispatch occurred.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`
- `benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt`
