---
id: T03
parent: S02
milestone: M022-i079tk
key_files:
  - .github/workflows/go-quality.yml
  - Dockerfile.onnx
  - tools/build_onnx_image.sh
  - docs/onnx-artifacts/README.md
key_decisions:
  - CI-safe checks are verified with actionlint and local equivalents.
  - Full ONNX image packaging remains locally proven but not hosted-CI-enabled until external artifacts are provisioned.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:48:08.392Z
blocker_discovered: false
---

# T03: Verified the CI boundary and full local ONNX Docker packaging proof.

**Verified the CI boundary and full local ONNX Docker packaging proof.**

## What Happened

Ran closure verification for S02 and M022. The GitHub workflow passed actionlint. The CI-safe allow-missing verifier passed. Default Go tests, pinned lint, native tokenizer tests, and ONNX+native smoke tests passed. Default Docker build passed. The dedicated ONNX image rebuilt and the container served `/health` plus a 1024-dimensional embedding smoke response. Binary hygiene, background cleanup, port cleanup, and GitNexus scope checks all passed.

## Verification

All closure checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/go-quality.yml` | 0 | ✅ pass — actionlint no findings | 5100ms |
| 2 | `python3 -m py_compile tools/verify_onnx_artifacts.py && tools/verify_onnx_artifacts.py --allow-missing ...` | 0 | ✅ pass — ci_allow_missing_verifier=pass | 5100ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass — Go test: 74 passed in 4 packages | 5100ms |
| 4 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 5000ms |
| 5 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — Go test: 16 passed in 1 package | 5000ms |
| 6 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — Go test: 2 passed in 1 package | 4900ms |
| 7 | `docker build -f api/Dockerfile -t fd-api:m022-default-final api` | 0 | ✅ pass — Successfully tagged fd-api:m022-default-final | 4800ms |
| 8 | `IMAGE_TAG=fd-api:onnx1024-m022-final tools/build_onnx_image.sh` | 0 | ✅ pass — Successfully tagged fd-api:onnx1024-m022-final | 53900ms |
| 9 | `curl -fsS http://localhost:18000/health && /v1/embeddings smoke` | 0 | ✅ pass — health ok; embedding_dims=1024 model=deepvk/USER-bge-m3 | 0ms |
| 10 | `binary hygiene, cleanup, GitNexus detect` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0; no background processes; port_18000_clean; GitNexus low scope | 0ms |

## Deviations

None.

## Known Issues

ONNX image build context is 1.569GB. Hosted CI should not run this until artifact provisioning/cache is designed.

## Files Created/Modified

- `.github/workflows/go-quality.yml`
- `Dockerfile.onnx`
- `tools/build_onnx_image.sh`
- `docs/onnx-artifacts/README.md`
