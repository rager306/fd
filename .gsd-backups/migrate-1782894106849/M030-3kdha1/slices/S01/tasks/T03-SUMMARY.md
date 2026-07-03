---
id: T03
parent: S01
milestone: M030-3kdha1
key_files:
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
  - api/main_test.go
  - tools/provision_onnx_artifacts.py
  - tools/verify_onnx_artifacts.py
  - tools/build_onnx_image.sh
key_decisions:
  - Pre-commit GitNexus HIGH scope is expected because the remediation touches Go manifest validation plus Python provisioning/verifier tooling. Direct impacts were checked before editing key affected symbols and all guardrails pass.
duration: 
verification_result: passed
completed_at: 2026-05-21T05:30:47.558Z
blocker_discovered: false
---

# T03: Verified M030 S01 path security remediation guardrails.

**Verified M030 S01 path security remediation guardrails.**

## What Happened

Ran S01 verification. Artifact tool guardrails passed, default Go tests passed with 87 tests after fixing the test fixture to use approved artifact roots, pinned lint had 0 issues, tagged tokenizer tests passed with 20 tests, ONNX smoke tests passed with 2 tests, actionlint passed, default Docker build passed, binary hygiene passed, port 18000 is clean, and no background processes remain. GitNexus reports expected high pre-commit scope across Go manifest and Python artifact tooling paths.

## Verification

All executable S01 checks passed; GitNexus pre-commit scope recorded as expected.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py benchmark.py tools/evaluate_legal_retrieval.py && provisioning dry-run/missing-source/verifier allow-missing` | 0 | ✅ pass — m030_artifact_tool_guardrails=pass | 8600ms |
| 2 | `cd api && go test ./... -short` | 0 | ✅ pass — Go test: 87 passed in 4 packages | 8500ms |
| 3 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 8500ms |
| 4 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — Go test: 20 passed in 1 package | 8400ms |
| 5 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — Go test: 2 passed in 1 package | 8400ms |
| 6 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/go-quality.yml .github/workflows/onnx-packaging.yml` | 0 | ✅ pass — no output | 8300ms |
| 7 | `docker build -f api/Dockerfile -t fd-api:m030-default-s01-final api` | 0 | ✅ pass — Successfully tagged fd-api:m030-default-s01-final | 8200ms |
| 8 | `binary hygiene, port cleanup, bg_shell list` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0; port_18000_clean; no background processes | 0ms |
| 9 | `gitnexus_detect_changes` | 0 | ⚠️ expected pre-commit HIGH scope — Go manifest and Python artifact tooling paths changed | 0ms |

## Deviations

Initial S01 full test run failed because `api/main_test.go` still generated an unapproved temporary absolute artifact path. The fixture was corrected to use the approved `.gsd/runtime/onnx/...` layout; rerun guardrails passed.

## Known Issues

Final post-commit reindex/detect required. Hosted workflow proof and immutable source selection remain future work.

## Files Created/Modified

- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`
- `api/main_test.go`
- `tools/provision_onnx_artifacts.py`
- `tools/verify_onnx_artifacts.py`
- `tools/build_onnx_image.sh`
