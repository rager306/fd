---
id: T03
parent: S02
milestone: M033-aym8ih
key_files:
  - tools/provision_onnx_artifacts.py
  - docs/onnx-artifacts/PROVISIONING.md
  - benchmark-results/fd-onnx-runtime-wheel-provisioning-m033-s02.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T07:24:07.787Z
blocker_discovered: false
---

# T03: Ran final M033 guardrails and confirmed ONNX Runtime wheel provisioning support is ready for closure.

**Ran final M033 guardrails and confirmed ONNX Runtime wheel provisioning support is ready for closure.**

## What Happened

Ran the full final guardrail set after the last code/docs/outcome changes. The provisioning helper compiles, dry-run and artifact/export verifiers pass, synthetic wheel probes cover positive extraction, missing member, symlink member rejection, and direct fallback, and project tests/lint/build checks pass. GitNexus detects HIGH scope because core provisioning functions changed, but affected flows are expected and verified.

## Verification

Fresh verification passed after final zip member regular-file handling change.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py tools/verify_onnx_export_contract.py && python3 tools/provision_onnx_artifacts.py --dry-run ... && python3 tools/verify_onnx_artifacts.py --allow-missing ... && python3 tools/verify_onnx_export_contract.py` | 0 | ✅ pass — compile/dry-run/verifiers passed | 25800ms |
| 2 | `M033 final synthetic wheel probes including symlink-like member` | 0 | ✅ pass — wheel_positive, missing_member, symlink_member, direct_fallback | 25800ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass — 87 passed in 4 packages | 25700ms |
| 4 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 25600ms |
| 5 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 25600ms |
| 6 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — 20 passed in 1 package | 9200ms |
| 7 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — 2 passed in 1 package | 9100ms |
| 8 | `docker build -f api/Dockerfile -t fd-api:m033-default-final api` | 0 | ✅ pass — image fd-api:m033-default-final built | 9000ms |
| 9 | `docs/outcome leak and signed URL check` | 0 | ✅ pass — missing_markers=0, raw_input_or_secret_leaks=0, signed_url_like=0 | 9100ms |
| 10 | `tracked binary hygiene via git ls-files` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0 | 9000ms |
| 11 | `gitnexus_detect_changes(scope=all)` | 0 | ✅ pass — expected HIGH scope in provisioning helper, 9 affected provisioning processes | 0ms |
| 12 | `bg_shell list and lsof port 18000` | 0 | ✅ pass — no background processes; port_18000_clean | 0ms |

## Deviations

Final synthetic probe initially failed after tightening the zip regular-file check because Python-created normal entries may not carry POSIX regular-file mode bits. The check was corrected to reject symlinks/directories/non-regular POSIX file types while tolerating entries with no POSIX file type, then all verification was rerun.

## Known Issues

GitNexus pre-commit risk is HIGH because this intentionally changes central provisioning helper flow (`main`/`materialize_source`). This is expected and covered by positive/negative synthetic probes plus project guardrails.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py`
- `docs/onnx-artifacts/PROVISIONING.md`
- `benchmark-results/fd-onnx-runtime-wheel-provisioning-m033-s02.txt`
