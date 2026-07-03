---
id: T02
parent: S02
milestone: M037-d23oz4
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/README.md
  - benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T10:21:47.698Z
blocker_discovered: false
---

# T02: Ran final M037 guardrails successfully.

**Ran final M037 guardrails successfully.**

## What Happened

Ran final M037 guardrails. Manifest JSON, Python tool compilation, provisioning dry-run, artifact verifier, export verifier, workflow syntax, Go tests, lint, tagged tests, Docker default build, target-runtime contract/leak checks, binary hygiene, background/port checks, and GitNexus detect all passed.

## Verification

All final M037 checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `json.tool + py_compile + provisioning dry-run + verify_onnx_artifacts + verify_onnx_export_contract` | 0 | ✅ pass — manifest valid and provisioning/export checks passed | 21400ms |
| 2 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 21300ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass — 87 passed in 4 packages | 21300ms |
| 4 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 21200ms |
| 5 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — 20 passed in 1 package | 8800ms |
| 6 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — 2 passed in 1 package | 8800ms |
| 7 | `docker build -f api/Dockerfile -t fd-api:m037-default-final api` | 0 | ✅ pass — image fd-api:m037-default-final built | 8700ms |
| 8 | `custom final target-runtime contract/leak checks` | 0 | ✅ pass — missing_markers=0, contract_checks_failed=0, raw_input_or_secret_leaks=0, signed_url_like=0 | 21200ms |
| 9 | `tracked binary hygiene` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0 | 8700ms |
| 10 | `gitnexus_detect_changes(scope=all)` | 0 | ✅ pass — low risk, no affected processes | 0ms |
| 11 | `bg_shell list and lsof port 18000` | 0 | ✅ pass — no background processes; port_18000_clean | 0ms |

## Deviations

None.

## Known Issues

Commit/reindex still pending after slice and milestone closure. No new target-runtime acceptance run was performed beyond existing guardrails.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/README.md`
- `benchmark-results/fd-onnx-target-runtime-validation-contract-m037-s01.txt`
