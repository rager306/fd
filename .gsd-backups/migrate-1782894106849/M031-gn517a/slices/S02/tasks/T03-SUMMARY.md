---
id: T03
parent: S02
milestone: M031-gn517a
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - docs/onnx-artifacts/PROVISIONING.md
  - benchmark-results/fd-onnx-source-contract-m031-s02.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T06:40:59.252Z
blocker_discovered: false
---

# T03: Ran final M031 guardrails and confirmed the source-contract changes are ready for closure.

**Ran final M031 guardrails and confirmed the source-contract changes are ready for closure.**

## What Happened

Ran final verification after the docs/manifest changes. Corrected two verification-command mistakes and reran the affected checks. GitNexus detects only low-risk documentation-section changes and no affected processes. Background process and port checks are clean.

## Verification

Fresh verification passed after the final source contract edits.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py && python3 tools/provision_onnx_artifacts.py --dry-run ... && python3 tools/verify_onnx_artifacts.py --allow-missing ...` | 0 | ✅ pass — dry-run shows expected blocked/missing sources and verifier validates local artifacts | 4400ms |
| 2 | `cd api && go test ./... -short` | 0 | ✅ pass — 87 passed in 4 packages | 29400ms |
| 3 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 29300ms |
| 4 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 29200ms |
| 5 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — 20 passed in 1 package | 14600ms |
| 6 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — 2 passed in 1 package | 14500ms |
| 7 | `docker build -f api/Dockerfile -t fd-api:m031-default-final api` | 0 | ✅ pass — image fd-api:m031-default-final built | 14400ms |
| 8 | `docs/outcome leak and signed URL check` | 0 | ✅ pass — missing_markers=0, raw_input_or_secret_leaks=0, signed_url_like=0 | 14500ms |
| 9 | `tracked binary hygiene via git ls-files` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0 | 4200ms |
| 10 | `gitnexus_detect_changes(scope=all)` | 0 | ✅ pass — low risk docs sections only; no affected processes | 0ms |
| 11 | `bg_shell list and lsof port 18000` | 0 | ✅ pass — no background processes; port_18000_clean | 0ms |

## Deviations

The first provisioning guardrail used a nonexistent `--allow-missing` flag for the provisioning helper, and the first binary hygiene check scanned ignored runtime artifacts instead of tracked files. Both were command-scope errors, corrected and rerun successfully.

## Known Issues

Local commit is still pending until after GSD slice/milestone closure and DB checkpoint.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `benchmark-results/fd-onnx-source-contract-m031-s02.txt`
