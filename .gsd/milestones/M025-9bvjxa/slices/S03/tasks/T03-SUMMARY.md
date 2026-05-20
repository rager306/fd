---
id: T03
parent: S03
milestone: M025-9bvjxa
key_files:
  - .github/workflows/onnx-packaging.yml
  - tools/provision_onnx_artifacts.py
  - docs/onnx-artifacts/PROVISIONING.md
  - docs/onnx-artifacts/OPERATIONS.md
  - docs/onnx-artifacts/README.md
key_decisions:
  - Manual ONNX workflow skeleton is valid but non-required; production path remains gated on artifact sources and operational proof.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:48:15.803Z
blocker_discovered: false
---

# T03: Verified the hosted ONNX CI skeleton and M025 guardrails.

**Verified the hosted ONNX CI skeleton and M025 guardrails.**

## What Happened

Ran S03 and milestone closure checks. Both workflows passed actionlint, scripts compiled, provisioning dry-run and allow-missing verifier passed, default Go tests passed, GolangCI-Lint had 0 issues, tagged tokenizer and ONNX smoke tests passed, default Docker build passed, binary hygiene passed, port 18000 is clean, no background processes remain, and GitNexus scope is low with docs-only indexed changes.

## Verification

All closure checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/go-quality.yml .github/workflows/onnx-packaging.yml` | 0 | ✅ pass — actionlint no findings | 9400ms |
| 2 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py tools/evaluate_legal_retrieval.py benchmark.py && provisioning dry-run + verifier allow-missing` | 0 | ✅ pass — m025_scripts_provision_verify=pass | 9300ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass — Go test: 74 passed in 4 packages | 9300ms |
| 4 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 9200ms |
| 5 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — Go test: 16 passed in 1 package | 9200ms |
| 6 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — Go test: 2 passed in 1 package | 9100ms |
| 7 | `docker build -f api/Dockerfile -t fd-api:m025-default-final api` | 0 | ✅ pass — Successfully tagged fd-api:m025-default-final | 9100ms |
| 8 | `binary hygiene, port cleanup, bg_shell list, GitNexus detect` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0; port_18000_clean; no background processes; GitNexus low scope | 0ms |

## Deviations

None.

## Known Issues

Full hosted workflow has not been run because real immutable artifact sources are not configured. This is expected and documented.

## Files Created/Modified

- `.github/workflows/onnx-packaging.yml`
- `tools/provision_onnx_artifacts.py`
- `docs/onnx-artifacts/PROVISIONING.md`
- `docs/onnx-artifacts/OPERATIONS.md`
- `docs/onnx-artifacts/README.md`
