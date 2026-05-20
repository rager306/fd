---
id: T03
parent: S01
milestone: M027-qswsja
key_files:
  - api/embed/onnx_manifest.go
  - api/embed/onnx_manifest_test.go
  - api/main.go
  - api/main_test.go
  - api/handlers/health.go
  - api/handlers/health_test.go
key_decisions:
  - S01 accepts pre-commit GitNexus high scope as expected implementation breadth; post-commit reindex/detect remains required.
duration: 
verification_result: passed
completed_at: 2026-05-20T12:39:57.935Z
blocker_discovered: false
---

# T03: Verified M027 S01 startup preflight diagnostics guardrails.

**Verified M027 S01 startup preflight diagnostics guardrails.**

## What Happened

Ran S01 verification. Default Go tests passed with 85 tests, pinned GolangCI-Lint had 0 issues, tagged native tokenizer tests passed with 18 tests, tagged ONNX smoke tests passed, default Docker build passed, actionlint passed, scripts/verifier passed, binary hygiene passed, port 18000 was clean, and no background processes remain. GitNexus reports expected high pre-commit scope across manifest/startup/health paths.

## Verification

All executable S01 guardrails passed; GitNexus high pre-commit scope recorded as expected.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./... -short` | 0 | ✅ pass — Go test: 85 passed in 4 packages | 39600ms |
| 2 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 39500ms |
| 3 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — Go test: 18 passed in 1 package | 39500ms |
| 4 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — Go test: 2 passed in 1 package | 39400ms |
| 5 | `docker build -f api/Dockerfile -t fd-api:m027-default-s01 api` | 0 | ✅ pass — Successfully tagged fd-api:m027-default-s01 | 39300ms |
| 6 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/go-quality.yml .github/workflows/onnx-packaging.yml` | 0 | ✅ pass — no output | 39200ms |
| 7 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py tools/evaluate_legal_retrieval.py benchmark.py && tools/verify_onnx_artifacts.py --allow-missing ...` | 0 | ✅ pass — m027_scripts_and_verifier=pass | 39200ms |
| 8 | `binary hygiene, port cleanup, bg_shell list` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0; port_18000_clean; no background processes | 0ms |
| 9 | `gitnexus_impact RuntimeHealth` | 0 | ✅ pass — LOW risk | 0ms |
| 10 | `gitnexus_detect_changes` | 0 | ⚠️ expected pre-commit high scope — manifest/startup/health/test paths | 0ms |

## Deviations

GitNexus pre-commit scope is high because startup config, manifest validation, health metadata, and tests changed. Direct impact checks were run for loadEmbeddingRuntimeConfig, ONNXArtifactManifest, ONNXArtifactValidation, and RuntimeHealth; risk is expected and covered by guardrails.

## Known Issues

Provider diagnostics are configured-provider validation only; runtime provider enumeration remains unavailable. Existing artifact validation errors may include filesystem paths, to be reviewed in a future security/logging gate.

## Files Created/Modified

- `api/embed/onnx_manifest.go`
- `api/embed/onnx_manifest_test.go`
- `api/main.go`
- `api/main_test.go`
- `api/handlers/health.go`
- `api/handlers/health_test.go`
