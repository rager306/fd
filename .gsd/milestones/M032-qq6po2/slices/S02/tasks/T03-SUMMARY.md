---
id: T03
parent: S02
milestone: M032-qq6po2
key_files:
  - tools/verify_onnx_export_contract.py
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - benchmark-results/fd-onnx-export-contract-verifier-m032-s02.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T07:02:11.739Z
blocker_discovered: false
---

# T03: Ran final M032 guardrails and confirmed the verifier/source strategy changes are ready for closure.

**Ran final M032 guardrails and confirmed the verifier/source strategy changes are ready for closure.**

## What Happened

Ran the full final guardrail set after the last documentation and decision changes. The new verifier passes positive checks, catches negative tamper probes, and project-level tests/lint/build checks pass. GitNexus reports low-risk docs changes only and no affected processes. Background process and port cleanup checks are clean.

## Verification

Fresh verification passed after the final changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/verify_onnx_export_contract.py tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py && python3 tools/verify_onnx_export_contract.py && python3 tools/provision_onnx_artifacts.py --dry-run ... && python3 tools/verify_onnx_artifacts.py --allow-missing ...` | 0 | ✅ pass — verifier passed; provisioning dry-run and local artifact verifier passed | 36200ms |
| 2 | `M032 final negative verifier probes` | 0 | ✅ pass — bad_sha, bad_revision, bad_transformers all failed with expected labels | 8300ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass — 87 passed in 4 packages | 36100ms |
| 4 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 36100ms |
| 5 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 36000ms |
| 6 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — 20 passed in 1 package | 8500ms |
| 7 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — 2 passed in 1 package | 8400ms |
| 8 | `docker build -f api/Dockerfile -t fd-api:m032-default-final api` | 0 | ✅ pass — image fd-api:m032-default-final built | 8300ms |
| 9 | `docs/outcome leak and signed URL check` | 0 | ✅ pass — missing_markers=0, raw_input_or_secret_leaks=0, signed_url_like=0 | 8400ms |
| 10 | `tracked binary hygiene via git ls-files` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0 | 36000ms |
| 11 | `gitnexus_detect_changes(scope=all)` | 0 | ✅ pass — low risk docs sections only; no affected processes | 0ms |
| 12 | `bg_shell list and lsof port 18000` | 0 | ✅ pass — no background processes; port_18000_clean | 0ms |

## Deviations

None. Removed generated `tools/__pycache__` before closure.

## Known Issues

Local commit is still pending until after milestone closure and DB checkpoint.

## Files Created/Modified

- `tools/verify_onnx_export_contract.py`
- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `benchmark-results/fd-onnx-export-contract-verifier-m032-s02.txt`
