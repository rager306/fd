---
id: T03
parent: S02
milestone: M034-9tfz77
key_files:
  - .github/workflows/onnx-packaging.yml
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-workflow-input-alignment-m034-s02.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T07:59:32.483Z
blocker_discovered: false
---

# T03: Ran final M034 guardrails and confirmed workflow input alignment is ready for closure.

**Ran final M034 guardrails and confirmed workflow input alignment is ready for closure.**

## What Happened

Ran final guardrails after workflow/docs/outcome/decision changes. Actionlint passes, provisioning/export checks pass, Go tests/lint/tagged tests pass, Docker default build passes, docs leak checks pass, tracked binary hygiene passes, and GitNexus reports low-risk docs/workflow changes with no affected processes.

## Verification

Fresh verification passed after final changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py tools/verify_onnx_export_contract.py && python3 tools/provision_onnx_artifacts.py --dry-run ... && python3 tools/verify_onnx_artifacts.py --allow-missing ... && python3 tools/verify_onnx_export_contract.py` | 0 | ✅ pass — compile/dry-run/verifiers passed | 29000ms |
| 2 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 28900ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass — 87 passed in 4 packages | 28900ms |
| 4 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 28800ms |
| 5 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — 20 passed in 1 package | 12700ms |
| 6 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — 2 passed in 1 package | 12600ms |
| 7 | `docker build -f api/Dockerfile -t fd-api:m034-default-final api` | 0 | ✅ pass — image fd-api:m034-default-final built | 12500ms |
| 8 | `docs/outcome leak and signed URL check` | 0 | ✅ pass — missing_markers=0, raw_input_or_secret_leaks=0, signed_url_like=0 | 28800ms |
| 9 | `tracked binary hygiene via git ls-files` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0 | 12600ms |
| 10 | `gitnexus_detect_changes(scope=all)` | 0 | ✅ pass — low risk docs/workflow changes, no affected processes | 0ms |
| 11 | `bg_shell list and lsof port 18000` | 0 | ✅ pass — no background processes; port_18000_clean | 0ms |

## Deviations

None.

## Known Issues

Local commit is still pending until after milestone closure and DB checkpoint. Exact ONNX binary source remains blocked.

## Files Created/Modified

- `.github/workflows/onnx-packaging.yml`
- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`
- `benchmark-results/fd-onnx-workflow-input-alignment-m034-s02.txt`
